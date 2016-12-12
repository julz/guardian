package runrunc

import (
	"os/exec"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/lager"

	"github.com/cloudfoundry/gunk/command_runner"
)

const DefaultRootPath = "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
const DefaultPath = "PATH=/usr/local/bin:/usr/bin:/bin"

// da doo
type RunRunc struct {
	commandRunner command_runner.CommandRunner
	runc          RuncBinary

	*Execer
	*DadooCreateStarter
	*OomWatcher
	*Statser
	*Stater
	*Killer
	*Deleter
}

//go:generate counterfeiter . RuncBinary
type RuncBinary interface {
	ExecCommand(id, processJSONPath, pidFilePath string) *exec.Cmd
	EventsCommand(id string) *exec.Cmd
	StateCommand(id, logFile string) *exec.Cmd
	StatsCommand(id, logFile string) *exec.Cmd
	KillCommand(id, signal, logFile string) *exec.Cmd
	DeleteCommand(id, logFile string) *exec.Cmd
}

func New(runner command_runner.CommandRunner, runcCmdRunner RuncCmdRunner, runc RuncBinary, dadooPath, runcPath string, execPreparer ExecPreparer, execRunner ExecRunner) *RunRunc {
	return &RunRunc{
		DadooCreateStarter: &DadooCreateStarter{execRunner, runner},
		Execer:             NewExecer(execPreparer, execRunner),

		OomWatcher: NewOomWatcher(runner, runc),
		Statser:    NewStatser(runcCmdRunner, runc),
		Stater:     NewStater(runcCmdRunner, runc),
		Killer:     NewKiller(runcCmdRunner, runc),
		Deleter:    NewDeleter(runcCmdRunner, runc),
	}
}

type DadooCreateStarter struct {
	dadooRunner   ExecRunner
	commandRunner command_runner.CommandRunner
}

func (r *DadooCreateStarter) Create(log lager.Logger, path, id string, io garden.ProcessIO) error {
	_, err := r.dadooRunner.Run(log, nil, path, id, nil, garden.ProcessIO{})
	return err
}

func (r *DadooCreateStarter) Start(log lager.Logger, id, path string) error {
	return r.commandRunner.Run(exec.Command("runc", "start", id))
}
