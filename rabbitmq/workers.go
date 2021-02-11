package rabbitmq

import (
    "fmt"
    "context"

    "go.uber.org/zap"
)

// workerManager will start go routines that read Rabbitmq messages.
func (g *QueueInput) workerManager(ctx context.Context) {
    for i := 0; i < g.rabbit.workerCount; i++ {
        go g.startWorker(ctx, i)
    }
}


// worker polls Rabbitmq and wrires messages to the pipeline.
func (g *QueueInput) startWorker(ctx context.Context, id int) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
        }

        g.rabbit.connecting.Wait()
        if err := g.worker(ctx, id); err != nil {
            g.Errorw("Worker error", zap.Error(err))
        }
    }
}

func  (g *QueueInput) worker(ctx context.Context, id int) error {
    g.Infow(fmt.Sprintf("Starting Rabbitmq consumer %d", id))
    for d := range g.rabbit.deliveries {
        g.rabbit.wg.Add(1)
        g.Write(ctx, parseMessage(d))
        // TODO: Ideally the message is acknowledged after the  output operator
        // successfully writes the entry to the log destination
        if err := d.Ack(false); err != nil {
            return err
        }
        g.rabbit.wg.Done()

        // check for closed context after writing current message, avoid breaking
        // this loop after taking a message off the channel
        select {
        case <-ctx.Done():
            g.Infow(fmt.Sprintf("Stopping Rabbitmq consumer %d, received shutdown", id))
            return nil
        default:
        }
    }
    g.Errorw(fmt.Sprintf("Stopping Rabbitmq consumer %d, channel closed", id))
    return nil
}
