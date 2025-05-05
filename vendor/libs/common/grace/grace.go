package grace

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// Connections #graceful коннекты должны уметь лишь одно - закрываться методом Close()
type Connections interface {
	Close()
}

// Service #graceful сервис должен показать имя методом Name() и уметь завершать работу методом Stop()
type Service interface {
	Name() string
	Stop()
}

type Logger interface {
	Printf(format string, args ...interface{})
}

type Graceful struct {
	fifo []Service
}

func KillThemSoftly(fifo ...Service) Graceful {
	return Graceful{fifo: fifo}
}

// Shutdown #graceful этой функцией, по идее, должно завершаться любое приложение
func (g Graceful) Shutdown(errs chan error, log Logger, connections Connections) {
	exitCode := 1

	go func() {
		c := make(chan os.Signal, 3)

		signal.Notify(c, syscall.SIGTERM)
		signal.Notify(c, syscall.SIGINT)
		signal.Notify(c, syscall.SIGILL)

		exitCode = 0

		errs <- fmt.Errorf("%s", <-c)
	}()

	// #graceful в этот канал прилетит либо сигнал ОС, либо критическая ошибка из какого-то сервиса
	err := <-errs
	// #graceful и тогда мы придем сюда, чтобы всех убить.
	log.Printf("terminated: %v\n", err)

	for _, svc := range g.fifo {
		if svc != nil {
			log.Printf("shutdown service %s started\n", svc.Name())
			svc.Stop()
			log.Printf("shutdown service %s completed\n", svc.Name())
		}
	}

	// #graceful не забываем закрыть соединения
	if connections != nil {
		connections.Close()
	}

	log.Printf("%s\n%s\n", "Connections closed", "Bye!")

	os.Exit(exitCode)
}

func NewService(name string, stop func()) Service {
	return service{
		name: name,
		stop: stop,
	}
}

// #graceful дефолтная реализация интерфейса, который умеет в graceful shutdown
type service struct {
	name string
	stop func()
}

func (s service) Stop() {
	s.stop()
}

func (s service) Name() string {
	return s.name
}
