# Task 10: Event Bus for Decoupling

**Priority:** Low
**Status:** Not Started
**Estimated Effort:** Medium
**Dependencies:** None (advanced feature)

## Problem Statement

Components are currently coupled through direct message passing:

- `OpenPanelMsg` directly couples panels to main model
- Hard to add observers or plugins
- Difficult to log or debug user actions
- Can't easily add analytics or telemetry
- Limited extensibility for plugins

This is acceptable for current scope, but limits future extensibility.

### Affected Files
- `tui/model.go` - Direct message handling
- `tui/components/panels.go` - Direct message creation
- Future plugin system

## Proposed Solution

Introduce an event bus for publish-subscribe pattern:

### 1. Event System

```go
// tui/events/event.go
package events

import "time"

// EventType identifies the type of event
type EventType string

const (
    PanelOpened        EventType = "panel.opened"
    PanelClosed        EventType = "panel.closed"
    NavigationForward  EventType = "navigation.forward"
    NavigationBackward EventType = "navigation.backward"
    TypeChanged        EventType = "type.changed"
    OverlayOpened      EventType = "overlay.opened"
    OverlayClosed      EventType = "overlay.closed"
    ItemSelected       EventType = "item.selected"
    FilterApplied      EventType = "filter.applied"
    ErrorOccurred      EventType = "error.occurred"
)

// Event represents something that happened in the application
type Event interface {
    // Type returns the event type
    Type() EventType

    // Timestamp returns when the event occurred
    Timestamp() time.Time

    // Data returns event-specific data
    Data() interface{}
}

// BaseEvent provides common event functionality
type BaseEvent struct {
    type      EventType
    timestamp time.Time
    data      interface{}
}

func NewBaseEvent(eventType EventType, data interface{}) BaseEvent {
    return BaseEvent{
        type:      eventType,
        timestamp: time.Now(),
        data:      data,
    }
}

func (e BaseEvent) Type() EventType       { return e.type }
func (e BaseEvent) Timestamp() time.Time  { return e.timestamp }
func (e BaseEvent) Data() interface{}     { return e.data }
```

### 2. Specific Events

```go
// tui/events/panel_events.go
package events

import "github.com/tonysyu/gqlxp/tui/components"

// PanelOpenedEvent is published when a panel is opened
type PanelOpenedEvent struct {
    BaseEvent
    Panel components.Panel
}

func NewPanelOpenedEvent(panel components.Panel) *PanelOpenedEvent {
    return &PanelOpenedEvent{
        BaseEvent: NewBaseEvent(PanelOpened, panel),
        Panel:     panel,
    }
}

// PanelClosedEvent is published when a panel is closed
type PanelClosedEvent struct {
    BaseEvent
    Panel components.Panel
}

func NewPanelClosedEvent(panel components.Panel) *PanelClosedEvent {
    return &PanelClosedEvent{
        BaseEvent: NewBaseEvent(PanelClosed, panel),
        Panel:     panel,
    }
}

// tui/events/navigation_events.go
type NavigationEvent struct {
    BaseEvent
    Direction  string // "forward" or "backward"
    FromPanel  int    // Panel index before navigation
    ToPanel    int    // Panel index after navigation
}

func NewNavigationForwardEvent(from, to int) *NavigationEvent {
    return &NavigationEvent{
        BaseEvent: NewBaseEvent(NavigationForward, nil),
        Direction: "forward",
        FromPanel: from,
        ToPanel:   to,
    }
}

// tui/events/type_events.go
type TypeChangedEvent struct {
    BaseEvent
    OldType string
    NewType string
}

func NewTypeChangedEvent(oldType, newType string) *TypeChangedEvent {
    return &TypeChangedEvent{
        BaseEvent: NewBaseEvent(TypeChanged, nil),
        OldType:   oldType,
        NewType:   newType,
    }
}

// tui/events/item_events.go
type ItemSelectedEvent struct {
    BaseEvent
    ItemName string
    ItemType string
}

func NewItemSelectedEvent(name, itemType string) *ItemSelectedEvent {
    return &ItemSelectedEvent{
        BaseEvent: NewBaseEvent(ItemSelected, nil),
        ItemName:  name,
        ItemType:  itemType,
    }
}
```

### 3. Event Bus

```go
// tui/events/bus.go
package events

import (
    "context"
    "sync"
)

// Handler is a function that handles an event
type Handler func(Event)

// EventBus manages event subscriptions and publishing
type EventBus struct {
    mu          sync.RWMutex
    subscribers map[EventType][]Handler
    asyncBus    chan Event
    ctx         context.Context
    cancel      context.CancelFunc
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
    ctx, cancel := context.WithCancel(context.Background())
    bus := &EventBus{
        subscribers: make(map[EventType][]Handler),
        asyncBus:    make(chan Event, 100), // Buffered channel
        ctx:         ctx,
        cancel:      cancel,
    }

    // Start async event processor
    go bus.processAsyncEvents()

    return bus
}

// Subscribe registers a handler for an event type
func (b *EventBus) Subscribe(eventType EventType, handler Handler) {
    b.mu.Lock()
    defer b.mu.Unlock()

    b.subscribers[eventType] = append(b.subscribers[eventType], handler)
}

// SubscribeAll registers a handler for all event types
func (b *EventBus) SubscribeAll(handler Handler) {
    b.mu.Lock()
    defer b.mu.Unlock()

    // Add to a special "all events" subscription
    b.subscribers["*"] = append(b.subscribers["*"], handler)
}

// Publish publishes an event synchronously
func (b *EventBus) Publish(event Event) {
    b.mu.RLock()
    defer b.mu.RUnlock()

    // Call handlers for specific event type
    if handlers, ok := b.subscribers[event.Type()]; ok {
        for _, handler := range handlers {
            handler(event)
        }
    }

    // Call handlers subscribed to all events
    if allHandlers, ok := b.subscribers["*"]; ok {
        for _, handler := range allHandlers {
            handler(event)
        }
    }
}

// PublishAsync publishes an event asynchronously
func (b *EventBus) PublishAsync(event Event) {
    select {
    case b.asyncBus <- event:
    case <-b.ctx.Done():
        // Bus is shutting down
    default:
        // Channel full, log warning but don't block
    }
}

// processAsyncEvents processes events from async channel
func (b *EventBus) processAsyncEvents() {
    for {
        select {
        case event := <-b.asyncBus:
            b.Publish(event)
        case <-b.ctx.Done():
            return
        }
    }
}

// Unsubscribe removes a handler (requires handler comparison, tricky in Go)
// For simplicity, could use subscription IDs instead
func (b *EventBus) UnsubscribeAll() {
    b.mu.Lock()
    defer b.mu.Unlock()

    b.subscribers = make(map[EventType][]Handler)
}

// Close shuts down the event bus
func (b *EventBus) Close() {
    b.cancel()
    close(b.asyncBus)
}
```

### 4. Middleware/Plugins

Support plugins through event handlers:

```go
// tui/events/middleware.go
package events

// Middleware wraps event handlers with additional functionality
type Middleware func(Handler) Handler

// Logger middleware logs all events
func LoggerMiddleware() Middleware {
    return func(next Handler) Handler {
        return func(event Event) {
            log.Printf("[EVENT] %s at %s: %+v\n",
                event.Type(),
                event.Timestamp().Format(time.RFC3339),
                event.Data(),
            )
            next(event)
        }
    }
}

// Analytics middleware tracks event statistics
type AnalyticsMiddleware struct {
    counts map[EventType]int
    mu     sync.RWMutex
}

func NewAnalyticsMiddleware() *AnalyticsMiddleware {
    return &AnalyticsMiddleware{
        counts: make(map[EventType]int),
    }
}

func (m *AnalyticsMiddleware) Middleware() Middleware {
    return func(next Handler) Handler {
        return func(event Event) {
            m.mu.Lock()
            m.counts[event.Type()]++
            m.mu.Unlock()
            next(event)
        }
    }
}

func (m *AnalyticsMiddleware) GetCounts() map[EventType]int {
    m.mu.RLock()
    defer m.mu.RUnlock()

    counts := make(map[EventType]int, len(m.counts))
    for k, v := range m.counts {
        counts[k] = v
    }
    return counts
}

// Apply middleware to a handler
func ApplyMiddleware(handler Handler, middlewares ...Middleware) Handler {
    result := handler
    // Apply in reverse order so first middleware is outermost
    for i := len(middlewares) - 1; i >= 0; i-- {
        result = middlewares[i](result)
    }
    return result
}
```

### 5. Integration with Application

```go
// tui/model.go
type mainModel struct {
    // ... existing fields
    eventBus *events.EventBus
}

func newModel(schema adapters.SchemaView) mainModel {
    eventBus := events.NewEventBus()

    // Subscribe to events
    eventBus.Subscribe(events.PanelOpened, func(e events.Event) {
        // Handle panel opened
    })

    m := mainModel{
        // ... existing init
        eventBus: eventBus,
    }

    return m
}

func (m *mainModel) handleOpenPanel(newPanel components.Panel) {
    // Publish event
    m.eventBus.Publish(events.NewPanelOpenedEvent(newPanel))

    // Existing logic
    m.panelStack = m.panelStack[:m.stackPosition+1]
    m.panelStack = append(m.panelStack, newPanel)
    m.sizePanels()
}

func (m *mainModel) incrementGQLTypeIndex(offset int) {
    oldType := string(m.selectedGQLType)

    // ... existing logic to change type

    newType := string(m.selectedGQLType)

    // Publish event
    m.eventBus.PublishAsync(events.NewTypeChangedEvent(oldType, newType))

    m.resetAndLoadMainPanel()
    m.sizePanels()
}
```

## Benefits

1. **Decoupling**: Components don't need direct references
2. **Extensibility**: Easy to add plugins via event handlers
3. **Observability**: Can log, monitor all events
4. **Testing**: Can verify events are published
5. **Analytics**: Track user behavior patterns
6. **Debugging**: Event log helps debugging

## Implementation Steps

1. Create `tui/events/` package
2. Define `Event` interface and base types
3. Implement `EventBus`
4. Define specific event types
5. Add tests for event bus and events
6. Integrate event bus into application
7. Publish events at key points
8. Add example middleware (logger, analytics)
9. Run tests: `just test`
10. Update documentation

## Testing Strategy

```go
// tui/events/bus_test.go
func TestEventBus_PublishSubscribe(t *testing.T) {
    bus := NewEventBus()
    defer bus.Close()

    var received Event
    bus.Subscribe(PanelOpened, func(e Event) {
        received = e
    })

    event := NewPanelOpenedEvent(testPanel)
    bus.Publish(event)

    assert.Equal(t, event, received)
}

func TestEventBus_SubscribeAll(t *testing.T) {
    bus := NewEventBus()
    defer bus.Close()

    count := 0
    bus.SubscribeAll(func(e Event) {
        count++
    })

    bus.Publish(NewPanelOpenedEvent(testPanel))
    bus.Publish(NewNavigationForwardEvent(0, 1))

    assert.Equal(t, 2, count)
}

// tui/events/middleware_test.go
func TestAnalyticsMiddleware(t *testing.T) {
    analytics := NewAnalyticsMiddleware()
    handler := ApplyMiddleware(
        func(e Event) {},
        analytics.Middleware(),
    )

    handler(NewPanelOpenedEvent(testPanel))
    handler(NewPanelOpenedEvent(testPanel))

    counts := analytics.GetCounts()
    assert.Equal(t, 2, counts[PanelOpened])
}
```

## Potential Issues

- **Complexity**: Adds abstraction layer
- **Performance**: Event publishing has overhead
- **Debugging**: Harder to trace event flow
- **Memory**: Event bus holds references

## Future Enhancements

1. **Plugin system**: Load plugins that subscribe to events
2. **Event replay**: Record and replay event sequences
3. **Remote events**: Publish events to remote systems
4. **Event filtering**: Filter events by criteria
5. **Priority handlers**: Execute some handlers before others
6. **Event aggregation**: Combine multiple events

## Related Tasks

- **Task 05** (Command Pattern): Commands could publish events
- **Task 02** (Navigation Manager): Publish navigation events
- **Task 09** (Configuration): Configure event handling behavior

## Notes

This is a low priority task because:
- Current direct message passing works well
- Event bus adds complexity
- Main benefit is extensibility for plugins
- Only implement if plugin system becomes important
