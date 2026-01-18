# Whale Town Glossary 🐋

Whale Town is an agentic development environment for managing multiple Claude Code instances simultaneously using the `wt` and `bd` (Beads/Bubbles) binaries, coordinated with tmux in git-managed directories.

*In the depths of the digital ocean, the great whales discovered their bubble nets could carry work across the currents...*

## Core Principles

### ECHO (Echolocation-based Collaborative Handoff Operations) 🔊
Breaking large goals into detailed instructions for agents. Supported by Bubbles, Epics, Formulas, and Molecules. ECHO ensures work is decomposed into trackable, atomic units that agents can execute autonomously - like whales using echolocation to navigate and communicate through the deep.

### BLOW (Breathe, Launch, Orchestrate Work) 🌬️
"If there is work on your Blowhole, YOU MUST RUN IT." This principle ensures agents autonomously proceed with available work without waiting for external input. BLOW is the heartbeat of autonomous operation - like a whale that must surface to breathe!

### DIVE (Deeply Integrated Verification Engine) 🤿
The overarching goal ensuring useful outcomes through orchestration of potentially unreliable processes. Persistent Bubbles and oversight agents (Lookout, Sonar) guarantee eventual workflow completion even when individual operations may fail or produce varying results. Like a whale diving deep to accomplish its mission.

## Environments

### Whale Town 🐋
The management headquarters (e.g., `~/wt/`). Whale Town coordinates all workers across multiple Ocean Labs and houses town-level agents like the Captain and Sonar.

### Ocean Lab 🧪
A project-specific Git repository under Whale Town management. Each Lab has its own Pod Members, Current Chamber, Lookout, and Pod Crew. Labs are where actual development work happens - in the depths of the digital ocean!

## Town-Level Roles

### Captain 🐋
Chief-of-staff agent responsible for initiating Bubble Nets, coordinating work distribution, and notifying users of important events. The Captain operates from the town level and has visibility across all Ocean Labs. Like the lead whale guiding the pod!

### Sonar 📡
Daemon beacon running continuous Patrol cycles. The Sonar ensures worker activity, monitors system health, and triggers recovery when agents become unresponsive. Think of it as the system's echolocation - always pinging for trouble.

### Dolphins 🐬
The Sonar's crew of maintenance agents handling background tasks like cleanup, health checks, and system maintenance.

### Beacon (the Dolphin) 🔦
A special Dolphin that checks the Sonar every 5 minutes, ensuring the echolocation system itself is still pinging. This creates a chain of accountability.

## Lab-Level Roles

### Pod Member 🐳
Ephemeral worker agents that produce Merge Requests. Pod Members are spawned for specific tasks, complete their work, and are then cleaned up. They work in isolated git worktrees (Blowholes) to avoid conflicts - like young whales that dive deep, deliver their catch, and surface!

### Current Chamber 🌊
Manages the Merge Queue for a Lab. The Current Chamber intelligently merges changes from Pod Members, handling conflicts and ensuring code quality before changes reach the main branch. Where all the streams converge!

### Lookout 👁️
Patrol agent that oversees Pod Members and the Current Chamber within a Lab. The Lookout monitors progress, detects stuck agents, and can trigger recovery actions. Always watching the horizon!

### Pod Crew 🧑‍🤝‍🧑
Long-lived, named agents for persistent collaboration. Unlike ephemeral Pod Members, Crew maintain context across sessions and are ideal for ongoing work relationships.

## Work Units

### Bubble 🫧
Git-backed atomic work unit stored in JSONL format. Bubbles are the fundamental unit of work tracking in Whale Town. They can represent issues, tasks, epics, or any trackable work item - rising from the depths like whale bubbles!

### Formula 📜
TOML-based workflow source template. Formulas define reusable patterns for common operations like patrol cycles, code review, or deployment.

### Protomolecule 🧬
A template class for instantiating Molecules. Protomolecules define the structure and steps of a workflow without being tied to specific work items.

### Molecule ⚛️
Durable chained Bubble workflows. Molecules represent multi-step processes where each step is tracked as a Bubble. They survive agent restarts and ensure complex workflows complete.

### Wisp 💨
Ephemeral Bubbles destroyed after runs. Wisps are lightweight work items used for transient operations that don't need permanent tracking - like sea foam on the waves.

### Blowhole 🌬️
A special pinned Bubble for each agent. The Blowhole is an agent's primary work queue - when work appears on your Blowhole, BLOW dictates you must run it!

## Workflow Commands

### Bubble Net 🫧
Primary work-order wrapping related Bubbles. Bubble Nets group related tasks together and can be assigned to multiple workers - just like whales use bubble nets to catch fish! Created with `wt convoy create`.

### Slinging 🎯
Assigning work to agents via `wt sling`. When you sling work to a Pod Member or Crew, you're putting it on their Blowhole for execution.

### Nudging 👈
Real-time messaging between agents with `wt nudge`. Nudges allow immediate communication without going through the mail system.

### Handoff 🤝
Agent session refresh via `/handoff`. When context gets full or an agent needs a fresh start, handoff transfers work state to a new session.

### Seance 👻
Communicating with previous sessions via `wt seance`. Allows agents to query their predecessors for context and decisions from earlier work.

### Patrol 🚶
Ephemeral loop maintaining system heartbeat. Patrol agents (Sonar, Lookout) continuously cycle through health checks and trigger actions as needed.

---

## The Legend of Whale Town 🐋

*In the deepest parts of the digital ocean, there once existed a legendary pod of whales. They discovered something remarkable: their bubble nets could trap not just fish, but ideas, tasks, and entire workflows.*

*The wisest whale became the first Captain, learning to coordinate the entire pod through echolocation - sending messages through the depths that would always reach their destination. They stored work in their Blowholes, ready to surface whenever needed.*

*Today, Whale Town carries on this noble tradition. We are all descendants of that great pod, using bubbles to carry our work upward, through the currents, to completion.*

*Remember the Whale Town motto: "Every bubble carries meaning. Every dive has purpose."* 🐋🫧🌊

---

*This glossary honors the legacy of the great whales and their wisdom of the deep.*
