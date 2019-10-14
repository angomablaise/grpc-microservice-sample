package repository_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/go-redis/redis/v7"
	"github.com/google/go-cmp/cmp"
	"github.com/ory/dockertest"
	"github.com/smockoro/grpc-microservice-sample/pkg/api"
	repo "github.com/smockoro/grpc-microservice-sample/pkg/repository/redis/go-redis"
)

var resource *dockertest.Resource

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect docker: %s", err)
	}
	resource, err = pool.Run("redis", "latest", nil)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		db := redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("localhost:%s", resource.GetPort("6379/tcp")),
		})

		return db.Ping().Err()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()

	err = pool.Purge(resource)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}

func TestSet(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("localhost:%s", resource.GetPort("6379/tcp")),
		Password: "",
		DB:       0, // use dafault DB
	})

	rrepo := repo.NewStoreRepository(client)

	cases := []struct {
		name string
		f    func(t *testing.T)
	}{
		{name: "Set Success", f: func(t *testing.T) {
			t.Parallel()
			store := &api.Store{
				Id:      1,
				Name:    "sample",
				Mail:    "sample@a.com",
				Address: "Tokyo",
			}

			ctx := context.Background()
			_, err := rrepo.Set(ctx, store)
			if err != nil {
				t.Error("err is nil but, err is ", err)
			}

			getData, _ := rrepo.Get(ctx, strconv.FormatInt(store.Id, 10))
			if diff := cmp.Diff(store, getData); diff != "" {
				t.Errorf("Set Data is %v but, Get Data is %v", store, getData)
			}
		}},
		{name: "Set failue Bad Host", f: func(t *testing.T) {
			t.Parallel()
			client := redis.NewClient(&redis.Options{
				Addr: "badhost:6379",
			})

			store := &api.Store{}
			rrepo := repo.NewStoreRepository(client)
			ctx := context.Background()
			_, err := rrepo.Set(ctx, store)
			if err == nil {
				t.Error("err is not but, err is nil")
			}
		}},
	}

	for _, c := range cases {
		t.Run(c.name, c.f)
	}
}

func TestGet(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("localhost:%s", resource.GetPort("6379/tcp")),
		Password: "",
		DB:       0, // use dafault DB
	})

	rrepo := repo.NewStoreRepository(client)

	cases := []struct {
		name string
		f    func(t *testing.T)
	}{
		{name: "Get Success", f: func(t *testing.T) {
			t.Parallel()
			setData := &api.Store{
				Id:      1,
				Name:    "sample",
				Mail:    "sample@a.com",
				Address: "Tokyo",
			}
			ctx := context.Background()
			_, _ = rrepo.Set(ctx, setData)

			store, err := rrepo.Get(ctx, "1")
			if err != nil {
				t.Error("err is nil but, err is ", err)
			}
			if diff := cmp.Diff(store, setData); diff != "" {
				t.Errorf("Set Data is %v but, Get Data is %v", setData, store)
			}

		}},
		{name: "Get failue Bad Host", f: func(t *testing.T) {
			t.Parallel()
			client := redis.NewClient(&redis.Options{
				Addr: "badhost:6379",
			})

			rrepo := repo.NewStoreRepository(client)
			ctx := context.Background()
			_, err := rrepo.Get(ctx, "1")
			if err == nil {
				t.Error("err is not but, err is nil")
			}
		}},
	}

	for _, c := range cases {
		t.Run(c.name, c.f)
	}
}

func TestDel(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0, // use dafault DB
	})

	rrepo := repo.NewStoreRepository(client)

	cases := []struct {
		name string
		f    func(t *testing.T)
	}{
		{name: "Del Success", f: func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			setData := &api.Store{
				Id:      1,
				Name:    "sample",
				Mail:    "sample@a.com",
				Address: "Tokyo",
			}
			_, _ = rrepo.Set(ctx, setData)

			delkey, err := rrepo.Del(ctx, strconv.FormatInt(setData.Id, 10))
			if err != nil && delkey == strconv.FormatInt(setData.Id, 10) {
				t.Error("err is not nil but, err is ", err)
			}
			_, err = rrepo.Get(ctx, delkey)
			if err == nil {
				t.Error(" err is nil but expected err is not nil")
			}
		}},
		{name: "Del failue Bad Host", f: func(t *testing.T) {
			t.Parallel()
			client := redis.NewClient(&redis.Options{
				Addr: "badhost:6379",
			})

			rrepo := repo.NewStoreRepository(client)
			ctx := context.Background()
			_, err := rrepo.Del(ctx, "1")
			if err == nil {
				t.Error("err is not but, err is nil")
			}
		}},
	}

	for _, c := range cases {
		t.Run(c.name, c.f)
	}
}

func BenchmarkSet(b *testing.B) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0, // use dafault DB
	})

	rrepo := repo.NewStoreRepository(client)
	setData := &api.Store{
		Name:    "sample",
		Mail:    "sample@a.com",
		Address: "Tokyo",
	}
	for i := 0; i < b.N; i++ {
		setData.Id = int64(i + 1)
		ctx := context.Background()
		b.StartTimer()
		_, err := rrepo.Set(ctx, setData)
		b.StopTimer()
		if err != nil {
			b.Fatal("Error")
		}
	}
}
