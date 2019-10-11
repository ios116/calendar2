package cmd

import (
	"github.com/ios116/calendar/internal/grpcserver"
	"github.com/spf13/cobra"
	"log"
)

var GrpcServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Run grpc server",
	Run: func(cmd *cobra.Command, args []string) {
		server()
	},
}

func server() {

	container := BuildContainer()
	err := container.Invoke(func(serverGRPS *grpcserver.GRPCServer) {
		serverGRPS.Start()
	})
	if err != nil {
		log.Println(err)
	}

	//appConf := config.NewAppConf()
	//logger, err := config.CreateLogger(appConf)
	//if err != nil {
	//	log.Fatalf("logger :%s", err.Error())
	//}
	//
	//
	//
	//dbConf := config.NewDateBaseConf()
	//conn, err := config.CreateDB(dbConf)
	//if err != nil {
	//	logger.Fatal("database:", zap.String("error", err.Error()))
	//}
	//if err = conn.Ping(); err != nil {
	//	logger.Fatal("can't connect to data base", zap.String("error", err.Error()))
	//}
	//
	//ctx := context.Background()
	//repository := store.NewDbRepo(ctx, conn)
	//calendarUseCase := calendar.NewCalendar(ctx, repository, logger)
	//
	//grpConf := config.NewGrpcConf()
	//serverGRPS := grpcserver.NewGRPCServer(logger, calendarUseCase, grpConf)
	//serverGRPS.Start()

	//conn.Close()

}
