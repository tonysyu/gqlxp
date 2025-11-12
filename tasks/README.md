# Architectural Improvement Tasks

This directory contains detailed documentation for proposed architectural improvements to gqlxp.

## Priority Guide

**High Priority** - Addresses current issues, reduces coupling, improves maintainability
- Task 01: Type Resolver Service
- Task 02: Navigation State Manager
- Task 07: Error Handling Strategy

**Medium Priority** - Reduces duplication, improves extensibility
- Task 03: Type Registry
- Task 04: Consolidate Adapters with Generics
- Task 08: Separate View Models

**Low Priority** - Nice-to-have refinements and advanced features
- Task 05: Command Pattern
- Task 06: Panel Lifecycle Manager
- Task 09: Configuration Layer
- Task 10: Event Bus

## Implementation Order

Recommended implementation sequence considering dependencies:

1. **Task 01** - Type Resolver Service (foundation for other improvements)
2. **Task 07** - Error Handling Strategy (complements resolver)
3. **Task 02** - Navigation State Manager (reduces mainModel complexity)
4. **Task 04** - Consolidate Adapters (leverages resolver from Task 01)
5. **Task 08** - Separate View Models (builds on resolver)
6. **Task 03** - Type Registry (independent, improves extensibility)
7. **Task 06** - Panel Lifecycle Manager (refines panel management)
8. **Task 09** - Configuration Layer (independent refinement)
9. **Task 05** - Command Pattern (advanced feature)
10. **Task 10** - Event Bus (only if needed for extensibility)

## Task Files

- [task-01-type-resolver.md](task-01-type-resolver.md) - Type Resolver Service
- [task-02-navigation-manager.md](task-02-navigation-manager.md) - Navigation State Manager
- [task-03-type-registry.md](task-03-type-registry.md) - Type Registry
- [task-04-consolidate-adapters.md](task-04-consolidate-adapters.md) - Consolidate Adapters with Generics
- [task-05-command-pattern.md](task-05-command-pattern.md) - Command Pattern for Actions
- [task-06-panel-lifecycle.md](task-06-panel-lifecycle.md) - Panel Lifecycle Manager
- [task-07-error-handling.md](task-07-error-handling.md) - Error Handling Strategy
- [task-08-view-models.md](task-08-view-models.md) - Separate View Models
- [task-09-configuration.md](task-09-configuration.md) - Configuration Layer
- [task-10-event-bus.md](task-10-event-bus.md) - Event Bus
