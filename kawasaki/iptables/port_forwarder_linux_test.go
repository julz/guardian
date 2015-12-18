package iptables_test

import (
	"fmt"
	"net"
	"os/exec"

	"github.com/cloudfoundry-incubator/guardian/kawasaki"
	"github.com/cloudfoundry-incubator/guardian/kawasaki/iptables"
	"github.com/cloudfoundry/gunk/command_runner/linux_command_runner"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PortForwarder", func() {
	var (
		spec            *kawasaki.PortForwarderSpec
		forwarder       *iptables.PortForwarder
		containerHandle string
		externalPort    uint32
		containerPort   uint32
		iptablesChain   string
		externalIP      net.IP
		containerIP     net.IP
	)

	BeforeEach(func() {
		externalPort = uint32(1210 + GinkgoParallelNode())
		containerPort = uint32(2120 + GinkgoParallelNode())
		containerHandle = fmt.Sprintf("h-%d", GinkgoParallelNode())
		iptablesChain = ""
		externalIP = nil
		containerIP = nil
	})

	JustBeforeEach(func() {
		spec = &kawasaki.PortForwarderSpec{
			IPTableChain: iptablesChain,
			ExternalIP:   externalIP,
			ContainerIP:  containerIP,
			FromPort:     externalPort,
			ToPort:       containerPort,
		}

		forwarder = iptables.NewPortForwarder(linux_command_runner.New())
	})

	Context("when NetworkConfig is valid", func() {
		BeforeEach(func() {
			externalIP = net.ParseIP("127.0.0.1")
			containerIP = net.ParseIP("127.0.0.2")
			iptablesChain = fmt.Sprintf("chain-%s", containerHandle)

			createChainCmd := exec.Command("iptables", "-w", "-t", "nat", "-N", iptablesChain)
			Expect(createChainCmd.Run()).To(Succeed())
		})

		AfterEach(func() {
			// clean up rules created by PortForwarder
			deleteRuleCmd := exec.Command(
				"iptables", "-w", "-t", "nat",
				"-D", iptablesChain,
				"-d", externalIP.String(),
				"-p", "tcp",
				"-m", "tcp",
				"--dport", fmt.Sprint(externalPort),
				"-j", "DNAT",
				"--to-destination", fmt.Sprintf("%s:%d", containerIP.String(), containerPort),
			)
			Expect(deleteRuleCmd.Run()).To(Succeed())

			deleteChainCmd := exec.Command(
				"iptables", "-w", "-t", "nat",
				"-X", iptablesChain,
			)
			Expect(deleteChainCmd.Run()).To(Succeed())
		})

		It("creates an iptables chain and adds a rule to it", func() {
			err := forwarder.Forward(spec)
			Expect(err).NotTo(HaveOccurred())

			out, err := exec.Command("iptables", "--table", "nat", "-S").CombinedOutput()
			Expect(err).NotTo(HaveOccurred())

			ipTableRules := string(out)

			Expect(ipTableRules).To(ContainSubstring(iptablesChain))
			Expect(ipTableRules).To(ContainSubstring(containerHandle))
			Expect(ipTableRules).To(ContainSubstring(fmt.Sprint(externalPort)))
			Expect(ipTableRules).To(ContainSubstring(fmt.Sprint(containerPort)))
			Expect(ipTableRules).To(ContainSubstring(externalIP.String()))
			Expect(ipTableRules).To(ContainSubstring(containerIP.String()))
		})
	})

	Context("when NetworkConfig is invalid", func() {
		BeforeEach(func() {
			iptablesChain = fmt.Sprintf("chain-%s", containerHandle)
		})

		It("returns the error", func() {
			err := forwarder.Forward(spec)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Bad IP address"))
		})
	})

})