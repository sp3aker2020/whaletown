# 🐋 Whale Town - Example Scenarios

*Welcome to the depths! Here are some whale-themed scenarios to explore Whale Town.*

---

## Scenario 1: The Great Migration 🌊

*The pod must migrate to new waters. Captain Orca coordinates the journey.*

```bash
# Create a bubble net to track the migration
wt convoy create "Great Migration to v2.0" \
    wt-update-deps \
    wt-fix-tests \
    wt-update-docs

# The Captain delegates to pod members
wt sling wt-update-deps frontend-lab
wt sling wt-fix-tests backend-lab  
wt sling wt-update-docs docs-lab

# Monitor the pod's progress
wt convoy list
wt status
```

**Whale Lore**: *In the deep waters, the Captain uses echolocation to track each bubble. When all bubbles surface, the migration is complete.*

---

## Scenario 2: The Bubble Net Feast 🫧

*Multiple pod members work together to create a bubble net - trapping work like whales trap fish!*

```bash
# Create a collaborative bubble net
wt convoy create "Auth System Overhaul" \
    wt-jwt-tokens \
    wt-oauth-flow \
    wt-session-mgmt \
    wt-security-audit

# Spawn 4 pod members (workers) to handle each bubble
wt sling wt-jwt-tokens auth-lab --create
wt sling wt-oauth-flow auth-lab --create
wt sling wt-session-mgmt auth-lab --create
wt sling wt-security-audit auth-lab --create

# Watch the pod coordinate
wt dashboard --open
```

**Whale Lore**: *Humpback whales create bubble nets by swimming in circles while blowing bubbles. The fish are trapped inside. Similarly, our bubble nets trap related work items so nothing escapes!*

---

## Scenario 3: The Deep Dive 🤿

*A pod member must dive deep to solve a complex issue.*

```bash
# A critical bug surfaces!
wt sling wt-critical-bug backend-lab --create

# The pod member dives deep (works on the issue)
# ... time passes ...

# When done, signal completion
wt done wt-critical-bug "Fixed the memory leak in the service layer"

# The bubble surfaces! 
wt convoy check  # Auto-closes completed bubble nets
```

**Whale Lore**: *Sperm whales can dive 3,000 feet deep and hold their breath for 90 minutes. Our pod members dive deep into code, surfacing only when the work is complete.*

---

## Scenario 4: Pod Communication 📡

*Whales communicate through the depths using echolocation.*

```bash
# Captain sends a message to a pod member
wt nudge backend-lab/polecats/Moby "Check the memory usage on the auth service"

# Broadcast to the entire pod
wt broadcast "Stand down - Captain reviewing the migration plan"

# Pod member escalates an issue
wt escalate "Found a blocker in the OAuth flow - need human review"

# Check your mailbox (incoming bubbles)
wt mail inbox
```

**Whale Lore**: *Blue whale calls can be heard 1,000 miles away through the ocean. Our messaging system ensures no bubble - no matter how small - goes unheard.*

---

## Scenario 5: The Pod Hierarchy 🐋👑

*Understanding the pod structure.*

```
                    🐋 Captain (Mayor)
                    The lead whale, coordinates all
                         │
            ┌────────────┼────────────┐
            │            │            │
       🌊 Ocean Lab  🌊 Ocean Lab  🌊 Ocean Lab
        (frontend)    (backend)     (docs)
            │            │            │
       ┌────┴────┐  ┌────┴────┐  ┌────┴────┐
       🐳 🐳 🐳    🐳 🐳       🐳 🐳
       Pod Members (Polecats)
       Ephemeral workers
```

```bash
# Start the Captain (main coordinator)
wt mayor attach

# Inside the Captain session, spawn pod members
wt sling wt-issue-1 frontend-lab  # Spawns a pod member
wt sling wt-issue-2 backend-lab   # Spawns another

# See who's swimming
wt agents

# Check on a specific pod member
wt peek frontend-lab/polecats/Nemo
```

---

## Scenario 6: The Watchful Whale Eye 👁️

*Special whales keep watch over the pod.*

| Watcher | Role | Command |
|---------|------|---------|
| 📡 **Sonar** | Daemon that pings all workers | `wt deacon status` |
| 👁️ **Lookout** | Watches workers in each lab | `wt witness status` |
| 🐬 **Dolphins** | Helper agents for cleanup | `wt dog list` |
| 🌊 **Current Chamber** | Merges completed work | `wt refinery status` |

```bash
# Start all watchers
wt up

# Check system health
wt status
wt doctor

# View the dashboard
wt dashboard --port 8080
```

---

## Scenario 7: The Legend Continues 📖

*Create your own whale tale...*

```bash
# Initialize your ocean
wt install ~/whale-ocean --git
cd ~/whale-ocean

# Add an ocean lab (your Git repository)
wt rig add my-project https://github.com/you/repo.git

# Create your crew workspace (human workspace)
wt crew add captain --rig my-project

# Start your journey
wt mayor attach
# Tell the Captain what you want to build!
```

---

## 🌊 The Whale Town Motto

> *"Every bubble carries meaning. Every dive has purpose."*

**ECHO** - Echolocation-based Collaborative Handoff Operations 🔊
**BLOW** - Breathe, Launch, Orchestrate Work 🌬️  
**DIVE** - Deeply Integrated Verification Engine 🤿

---

*May your bubbles always surface!* 🐋🫧🌊
