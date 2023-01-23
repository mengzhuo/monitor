package mqserver

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/eclipse/paho.golang/paho"
	"go.opentelemetry.io/otel"
	otrace "go.opentelemetry.io/otel/trace"
)

type mqttResponseRouter struct {
	prefix string
	mu     sync.RWMutex
	rm     map[string]chan<- *paho.Publish
}

func (mqs *server) setupResponseRouter(ctx context.Context, topicPrefix string) *mqttResponseRouter {
	topicPrefix = strings.TrimSuffix(topicPrefix, "#")
	return &mqttResponseRouter{
		prefix: topicPrefix,
		rm:     make(map[string]chan<- *paho.Publish, 100),
	}
}

func (rr *mqttResponseRouter) Handler() paho.MessageHandler {

	prefixLength := len(rr.prefix)
	tp := otel.GetTracerProvider().Tracer("mqserver")

	return func(p *paho.Publish) {

		ctx := context.Background()

		var traceID otrace.TraceID
		var spanID otrace.SpanID

		if tr := p.Properties.User.Get("TraceID"); len(tr) > 0 {
			traceID, _ = otrace.TraceIDFromHex(tr)
			spanID, _ = otrace.SpanIDFromHex(p.Properties.User.Get("SpanID"))
		}

		spanContext := otrace.NewSpanContext(otrace.SpanContextConfig{
			TraceID: traceID,
			SpanID:  spanID,
			// TraceFlags: traceFlags,
		})

		ctx = otrace.ContextWithSpanContext(ctx, spanContext)

		ctx, span := tp.Start(ctx, "mqhandler")
		defer span.End()

		if deadline, ok := ctx.Deadline(); ok {
			log.Printf("has deadline: %s", deadline)
		}

		log.Printf("handling message: %+v", p)
		span.AddEvent(fmt.Sprintf("handling message: %+v", p))

		topic := p.Topic

		if len(topic) < prefixLength {
			log.Printf("message topic too short (%q)", topic)
			return
		}

		log.Printf("topic / prefix length: %q / %d", topic, prefixLength)

		topicPath := strings.Split(topic[prefixLength:], "/")

		if len(topicPath) < 2 {
			log.Printf("could not get host and id from topic: %s", topic)
			return
		}

		log.Printf("topic path: %+v", topicPath)

		rr.mu.RLock()
		defer rr.mu.RUnlock()

		if ch, ok := rr.rm[topicPath[1]]; ok {
			// todo: include host from topicPath[0]
			ch <- p
		} else {
			log.Printf("no response channel for %s", topicPath[1])
		}

	}
}

func (rr *mqttResponseRouter) AddResponseID(id string, rc chan<- *paho.Publish) {
	log.Printf("adding channel for %s", id)
	rr.mu.Lock()
	defer rr.mu.Unlock()
	rr.rm[id] = rc
}

func (rr *mqttResponseRouter) CloseResponseID(id string) {
	log.Printf("closing channel for %s", id)
	rr.mu.Lock()
	defer rr.mu.Unlock()
	close(rr.rm[id])
	delete(rr.rm, id)
}
