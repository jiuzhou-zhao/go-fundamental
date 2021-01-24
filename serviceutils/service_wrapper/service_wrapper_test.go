package service_wrapper

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jiuzhou-zhao/go-fundamental/interfaces"
	"github.com/stretchr/testify/assert"
)

type TestCycleJob struct {
	id int
}

func (job *TestCycleJob) DoJob(ctx context.Context, logger interfaces.Logger) (time.Duration, error) {
	fmt.Printf("[%v]id: %v\n", time.Now(), job.id)
	return time.Second * time.Duration(job.id), nil
}

func TestCycleServiceWrapper(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	sw := NewCycleServiceWrapper(ctx, nil)
	err = sw.Start(&TestCycleJob{id: 1})
	assert.Nil(t, err)
	err = sw.Start(&TestCycleJob{id: 2})
	assert.Nil(t, err)
	sw.Wait()
}
