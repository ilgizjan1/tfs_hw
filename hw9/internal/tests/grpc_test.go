package tests

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	grpchandler "homework9/internal/api/handlers/grpc"
	contracts "homework9/internal/api/handlers/grpc/contracts/langs/go"
	"homework9/internal/repository/local-repo"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"homework9/internal/service"
)

func TestGRRPCCreateUser(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	adService := service.NewAdService(localrepo.NewAdRepo())
	userService := service.NewUserService(localrepo.NewUserRepo())

	grpcAdHandler := grpchandler.NewAdHandler(adService)
	contracts.RegisterAdServiceServer(srv, grpcAdHandler)

	grpcUserHandler := grpchandler.NewUserHandler(userService)
	contracts.RegisterUserServiceServer(srv, grpcUserHandler)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(
		ctx,
		"",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})
	client := contracts.NewUserServiceClient(conn)
	res, err := client.CreateUser(ctx, &contracts.CreateUserRequest{Nickname: "Oleg", Email: "olega@gmail.com"})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, "Oleg", res.Nickname)
	assert.Equal(t, "olega@gmail.com", res.Email)
}

func TestGRRPCUpdateUser(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	adService := service.NewAdService(localrepo.NewAdRepo())
	userService := service.NewUserService(localrepo.NewUserRepo())

	grpcAdHandler := grpchandler.NewAdHandler(adService)
	contracts.RegisterAdServiceServer(srv, grpcAdHandler)

	grpcUserHandler := grpchandler.NewUserHandler(userService)
	contracts.RegisterUserServiceServer(srv, grpcUserHandler)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(
		ctx,
		"",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})
	client := contracts.NewUserServiceClient(conn)
	res, err := client.CreateUser(ctx, &contracts.CreateUserRequest{Nickname: "Oleg", Email: "olega@gmail.com"})
	assert.NoError(t, err, "client.GetUser")

	res, err = client.UpdateUser(ctx, &contracts.UpdateUserRequest{UserId: res.UserId, Nickname: "Alena", Email: "alena@gmail.com"})
	assert.NoError(t, err, "client.UpdateUser")

	assert.Equal(t, "Alena", res.Nickname)
	assert.Equal(t, "alena@gmail.com", res.Email)
}

func TestGRRPCGetUser(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	adService := service.NewAdService(localrepo.NewAdRepo())
	userService := service.NewUserService(localrepo.NewUserRepo())

	grpcAdHandler := grpchandler.NewAdHandler(adService)
	contracts.RegisterAdServiceServer(srv, grpcAdHandler)

	grpcUserHandler := grpchandler.NewUserHandler(userService)
	contracts.RegisterUserServiceServer(srv, grpcUserHandler)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(
		ctx,
		"",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})
	client := contracts.NewUserServiceClient(conn)
	res, err := client.CreateUser(ctx, &contracts.CreateUserRequest{Nickname: "Oleg", Email: "olega@gmail.com"})
	assert.NoError(t, err, "client.GetUser")

	res, err = client.GetUser(ctx, &contracts.GetUserRequest{UserId: res.UserId})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, "Oleg", res.Nickname)
	assert.Equal(t, "olega@gmail.com", res.Email)
}

func TestGRRPCDeleteUser(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	adService := service.NewAdService(localrepo.NewAdRepo())
	userService := service.NewUserService(localrepo.NewUserRepo())

	grpcAdHandler := grpchandler.NewAdHandler(adService)
	contracts.RegisterAdServiceServer(srv, grpcAdHandler)

	grpcUserHandler := grpchandler.NewUserHandler(userService)
	contracts.RegisterUserServiceServer(srv, grpcUserHandler)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(
		ctx,
		"",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})
	client := contracts.NewUserServiceClient(conn)
	res, err := client.CreateUser(ctx, &contracts.CreateUserRequest{Nickname: "Oleg", Email: "olega@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	_, err = client.DeleteUser(ctx, &contracts.DeleteUserRequest{UserId: res.UserId})
	assert.NoError(t, err, "client.DeleteUser")
}

func TestGRRPCGetAd(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	adService := service.NewAdService(localrepo.NewAdRepo())
	userService := service.NewUserService(localrepo.NewUserRepo())

	grpcAdHandler := grpchandler.NewAdHandler(adService)
	contracts.RegisterAdServiceServer(srv, grpcAdHandler)

	grpcUserHandler := grpchandler.NewUserHandler(userService)
	contracts.RegisterUserServiceServer(srv, grpcUserHandler)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(
		ctx,
		"",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})
	clientUser := contracts.NewUserServiceClient(conn)

	user, err := clientUser.CreateUser(ctx, &contracts.CreateUserRequest{Nickname: "Oleg", Email: "olega@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	clientAd := contracts.NewAdServiceClient(conn)

	res, err := clientAd.CreateAd(ctx, &contracts.CreateAdRequest{Title: "the book", Text: "the text", UserId: user.UserId})
	assert.NoError(t, err, "client.CreateAd")

	res, err = clientAd.GetAd(ctx, &contracts.GetAdRequest{AdId: res.Id})
	assert.NoError(t, err, "client.GetAd")

	assert.Equal(t, "the book", res.Title)
	assert.Equal(t, "the text", res.Text)
}

func TestGRRPCCreateAd(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	adService := service.NewAdService(localrepo.NewAdRepo())
	userService := service.NewUserService(localrepo.NewUserRepo())

	grpcAdHandler := grpchandler.NewAdHandler(adService)
	contracts.RegisterAdServiceServer(srv, grpcAdHandler)

	grpcUserHandler := grpchandler.NewUserHandler(userService)
	contracts.RegisterUserServiceServer(srv, grpcUserHandler)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(
		ctx,
		"",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})
	clientUser := contracts.NewUserServiceClient(conn)

	user, err := clientUser.CreateUser(ctx, &contracts.CreateUserRequest{Nickname: "Oleg", Email: "olega@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	clientAd := contracts.NewAdServiceClient(conn)

	res, err := clientAd.CreateAd(ctx, &contracts.CreateAdRequest{Title: "the book", Text: "the text", UserId: user.UserId})
	assert.NoError(t, err, "client.CreateAd")

	assert.Equal(t, "the book", res.Title)
	assert.Equal(t, "the text", res.Text)
}

func TestGRRPCUpdateAd(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	adService := service.NewAdService(localrepo.NewAdRepo())
	userService := service.NewUserService(localrepo.NewUserRepo())

	grpcAdHandler := grpchandler.NewAdHandler(adService)
	contracts.RegisterAdServiceServer(srv, grpcAdHandler)

	grpcUserHandler := grpchandler.NewUserHandler(userService)
	contracts.RegisterUserServiceServer(srv, grpcUserHandler)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(
		ctx,
		"",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})
	clientUser := contracts.NewUserServiceClient(conn)

	user, err := clientUser.CreateUser(ctx, &contracts.CreateUserRequest{Nickname: "Oleg", Email: "olega@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	clientAd := contracts.NewAdServiceClient(conn)

	res, err := clientAd.CreateAd(ctx, &contracts.CreateAdRequest{Title: "the book", Text: "the text", UserId: user.UserId})
	assert.NoError(t, err, "client.CreateAd")

	res, err = clientAd.UpdateAd(ctx, &contracts.UpdateAdRequest{AdId: res.Id,
		Title: "new book", Text: "new text", UserId: user.UserId})
	assert.NoError(t, err, "client.UpdateAd")

	assert.Equal(t, "new book", res.Title)
	assert.Equal(t, "new text", res.Text)
}

func TestGRRPCChangeAdStatus(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	adService := service.NewAdService(localrepo.NewAdRepo())
	userService := service.NewUserService(localrepo.NewUserRepo())

	grpcAdHandler := grpchandler.NewAdHandler(adService)
	contracts.RegisterAdServiceServer(srv, grpcAdHandler)

	grpcUserHandler := grpchandler.NewUserHandler(userService)
	contracts.RegisterUserServiceServer(srv, grpcUserHandler)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(
		ctx,
		"",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})
	clientUser := contracts.NewUserServiceClient(conn)

	user, err := clientUser.CreateUser(ctx, &contracts.CreateUserRequest{Nickname: "Oleg", Email: "olega@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	clientAd := contracts.NewAdServiceClient(conn)

	res, err := clientAd.CreateAd(ctx, &contracts.CreateAdRequest{Title: "the book", Text: "the text", UserId: user.UserId})
	assert.NoError(t, err, "client.CreateAd")

	res, err = clientAd.ChangeAdStatus(ctx, &contracts.ChangeAdStatusRequest{AdId: res.Id, UserId: user.UserId,
		Published: true})
	assert.NoError(t, err, "client.ChangeAdStatus")

	assert.Equal(t, true, res.Published)
}

func TestGRRPCDeleteAd(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	adService := service.NewAdService(localrepo.NewAdRepo())
	userService := service.NewUserService(localrepo.NewUserRepo())

	grpcAdHandler := grpchandler.NewAdHandler(adService)
	contracts.RegisterAdServiceServer(srv, grpcAdHandler)

	grpcUserHandler := grpchandler.NewUserHandler(userService)
	contracts.RegisterUserServiceServer(srv, grpcUserHandler)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(
		ctx,
		"",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})
	clientUser := contracts.NewUserServiceClient(conn)

	user, err := clientUser.CreateUser(ctx, &contracts.CreateUserRequest{Nickname: "Oleg", Email: "olega@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	clientAd := contracts.NewAdServiceClient(conn)

	res, err := clientAd.CreateAd(ctx, &contracts.CreateAdRequest{Title: "the book", Text: "the text", UserId: user.UserId})
	assert.NoError(t, err, "client.CreateAd")

	_, err = clientAd.DeleteAd(ctx, &contracts.DeleteAdRequest{AdId: res.Id})
	assert.NoError(t, err, "client.DeleteAd")
}

func TestGRRPCSearchAds(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	adService := service.NewAdService(localrepo.NewAdRepo())
	userService := service.NewUserService(localrepo.NewUserRepo())

	grpcAdHandler := grpchandler.NewAdHandler(adService)
	contracts.RegisterAdServiceServer(srv, grpcAdHandler)

	grpcUserHandler := grpchandler.NewUserHandler(userService)
	contracts.RegisterUserServiceServer(srv, grpcUserHandler)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(
		ctx,
		"",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})
	clientUser := contracts.NewUserServiceClient(conn)

	user, err := clientUser.CreateUser(ctx, &contracts.CreateUserRequest{Nickname: "Oleg", Email: "olega@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	clientAd := contracts.NewAdServiceClient(conn)

	_, err = clientAd.CreateAd(ctx, &contracts.CreateAdRequest{Title: "cat and dog", Text: "the text", UserId: user.UserId})
	assert.NoError(t, err, "client.CreateAd")

	res, err := clientAd.CreateAd(ctx, &contracts.CreateAdRequest{Title: "cats and dogs", Text: "the text", UserId: user.UserId})
	assert.NoError(t, err, "client.CreateAd")

	ads, err := clientAd.SearchAds(ctx, &contracts.SearchAdsRequest{Text: "cats"})
	assert.NoError(t, err, "client.SearchAds")

	assert.Len(t, ads.List, 1)
	assert.Equal(t, ads.List[0].Id, res.Id)
	assert.Equal(t, ads.List[0].Title, res.Title)
	assert.Equal(t, ads.List[0].Text, res.Text)
	assert.Equal(t, ads.List[0].UserId, res.UserId)
}

func TestGRRPCListAds(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	adService := service.NewAdService(localrepo.NewAdRepo())
	userService := service.NewUserService(localrepo.NewUserRepo())

	grpcAdHandler := grpchandler.NewAdHandler(adService)
	contracts.RegisterAdServiceServer(srv, grpcAdHandler)

	grpcUserHandler := grpchandler.NewUserHandler(userService)
	contracts.RegisterUserServiceServer(srv, grpcUserHandler)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(
		ctx,
		"",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})
	clientUser := contracts.NewUserServiceClient(conn)

	user, err := clientUser.CreateUser(ctx, &contracts.CreateUserRequest{Nickname: "Oleg", Email: "olega@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	clientAd := contracts.NewAdServiceClient(conn)

	_, err = clientAd.CreateAd(ctx, &contracts.CreateAdRequest{Title: "cat and dog", Text: "the text", UserId: user.UserId})
	assert.NoError(t, err, "client.CreateAd")

	res, err := clientAd.CreateAd(ctx, &contracts.CreateAdRequest{Title: "cats and dogs", Text: "the text", UserId: user.UserId})
	assert.NoError(t, err, "client.CreateAd")

	_, err = clientAd.ChangeAdStatus(ctx, &contracts.ChangeAdStatusRequest{AdId: res.Id, UserId: user.UserId, Published: true})
	assert.NoError(t, err, "client.ChangeAdStatus")

	ads, err := clientAd.ListAds(ctx, &contracts.ListAdsRequest{Published: "true", UserId: "0"})
	assert.NoError(t, err, "client.ListAds")

	assert.Len(t, ads.List, 1)
	assert.Equal(t, ads.List[0].Id, res.Id)
	assert.Equal(t, ads.List[0].Title, res.Title)
	assert.Equal(t, ads.List[0].Text, res.Text)
	assert.Equal(t, ads.List[0].UserId, res.UserId)
}
