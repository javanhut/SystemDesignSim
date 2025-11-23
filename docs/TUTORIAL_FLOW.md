# Tutorial System - User Flow Diagram

## Complete User Journey

```
┌─────────────────────────────────────────────────────────────┐
│                    APPLICATION LAUNCH                        │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
            ┌─────────────────────┐
            │ Load Preferences    │
            │ from ~/.config      │
            └─────────┬───────────┘
                      │
          ┌───────────┴───────────┐
          │                       │
    FirstLaunch?             Already Used?
          │                       │
          ▼                       ▼
┌──────────────────┐    ┌──────────────────┐
│ WELCOME SCREEN   │    │  LEVEL SELECT    │
│   (Tutorial)     │    │     SCREEN       │
└────────┬─────────┘    └─────────┬────────┘
         │                        │
         │                        │ [View Tutorial]
         │                        │      Button
         │                        ▼
         │              ┌──────────────────┐
         └─────────────▶│ WELCOME SCREEN   │
                        │   (Tutorial)     │
                        └────────┬─────────┘
                                 │
                                 ▼

┌────────────────────────────────────────────────────────────────┐
│                    TUTORIAL - 7 PAGES                           │
├────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Page 1: Welcome & Overview                                    │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │ • What is System Design Simulator                        │ │
│  │ • What you'll learn                                      │ │
│  │ • How the game works                                     │ │
│  │ • Educational goals                                      │ │
│  └──────────────────────────────────────────────────────────┘ │
│                         [Next →]                               │
│                                                                 │
│  Page 2: Understanding Scenarios                               │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │ Level 1: Local Blog (10 users)                          │ │
│  │   → Scenario: Blog for friends                          │ │
│  │   → Learn: Basic architecture                           │ │
│  │                                                          │ │
│  │ Level 2: Growing Blog (100 users)                       │ │
│  │   → Scenario: Viral blog                                │ │
│  │   → Learn: Load balancing                               │ │
│  │                                                          │ │
│  │ ... (all 5 levels with context)                         │ │
│  └──────────────────────────────────────────────────────────┘ │
│                  [← Previous] [Next →]                         │
│                                                                 │
│  Page 3: Components Part 1 (Compute & Storage)                │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │ API SERVER (Blue)                                        │ │
│  │   • Sizes: S/M/L/XL (10/50/200/500 concurrent)          │ │
│  │   • Processing: 10-15ms                                  │ │
│  │   • Cost: $0.05-$0.40/hour                              │ │
│  │   • When: Handle user requests                           │ │
│  │                                                          │ │
│  │ DATABASE (Purple)                                        │ │
│  │   • Types: SQL, NoSQL, Key-Value, Document              │ │
│  │   • Latency: 10ms read, 15ms write                      │ │
│  │   • Features: Sharding, Replication                     │ │
│  │   • When: Need persistent storage                       │ │
│  └──────────────────────────────────────────────────────────┘ │
│                  [← Previous] [Next →]                         │
│                                                                 │
│  Page 4: Components Part 2 (Performance)                      │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │ CACHE (Green)                                            │ │
│  │   • Policies: LRU, LFU, FIFO                            │ │
│  │   • Speed: 1-2ms (10x faster!)                          │ │
│  │   • Target: 70%+ hit rate                               │ │
│  │   • How: Check cache → DB on miss                       │ │
│  │                                                          │ │
│  │ LOAD BALANCER (Yellow)                                   │ │
│  │   • Strategies: Round-robin, Least-connected            │ │
│  │   • Benefits: HA, scalability                           │ │
│  │   • Pattern: LB → Multiple API Servers                  │ │
│  └──────────────────────────────────────────────────────────┘ │
│                  [← Previous] [Next →]                         │
│                                                                 │
│  Page 5: Components Part 3 (Global)                           │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │ CDN (Dark Blue)                                          │ │
│  │   • Edge Locations: 5 regions globally                   │ │
│  │   • Speed: 2ms edge vs 200ms cross-region               │ │
│  │   • 100x faster for global users!                       │ │
│  │                                                          │ │
│  │ NETWORK & LATENCY                                        │ │
│  │   • Regional latency matrix                              │ │
│  │   • Optimization tips                                    │ │
│  │   • Understanding P99                                    │ │
│  └──────────────────────────────────────────────────────────┘ │
│                  [← Previous] [Next →]                         │
│                                                                 │
│  Page 6: How to Play                                          │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │ ADDING COMPONENTS:                                       │ │
│  │   1. Click toolbox button                                │ │
│  │   2. Component appears on canvas                         │ │
│  │                                                          │ │
│  │ CONNECTING:                                              │ │
│  │   1. Right-click source                                  │ │
│  │   2. Drag to target                                      │ │
│  │   3. Release to connect                                  │ │
│  │   4. See animated particles!                            │ │
│  │                                                          │ │
│  │ HEALTH COLORS:                                           │ │
│  │   Green → Yellow → Orange → Red                         │ │
│  │                                                          │ │
│  │ RUNNING & SUBMITTING:                                    │ │
│  │   • Start → Monitor → Stop → Submit                     │ │
│  └──────────────────────────────────────────────────────────┘ │
│                  [← Previous] [Next →]                         │
│                                                                 │
│  Page 7: Scoring & Progression                                │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │ REQUIREMENTS (Must Pass):                                │ │
│  │   ✓ Max latency, Min uptime                             │ │
│  │   ✓ Budget constraint                                    │ │
│  │   ✓ Architecture requirements                            │ │
│  │                                                          │ │
│  │ BONUS POINTS:                                            │ │
│  │   • Excellent metrics: +100 each                         │ │
│  │   • Cost efficient: +150                                 │ │
│  │   • Cost savings: up to +200                            │ │
│  │                                                          │ │
│  │ FORMULA:                                                 │ │
│  │   Base: 1000                                             │ │
│  │   Penalties: -200 per failure                            │ │
│  │   Bonuses: up to +700                                    │ │
│  │   Total: 0-1700 points                                   │ │
│  │                                                          │ │
│  │ TIPS:                                                    │ │
│  │   • Start simple                                         │ │
│  │   • Use caching                                          │ │
│  │   • Monitor health colors                                │ │
│  │   • Balance cost vs performance                          │ │
│  └──────────────────────────────────────────────────────────┘ │
│                  [← Previous] [Get Started! →]                 │
│                                                                 │
│                                                                 │
│          [Skip Tutorial] (available on all pages)              │
│                                                                 │
└────────────────────────┬───────────────────────────────────────┘
                         │
                         │ Complete or Skip
                         ▼
                 ┌───────────────┐
                 │ Save Prefs    │
                 │ tutorial_     │
                 │ completed=true│
                 └───────┬───────┘
                         │
                         ▼
┌────────────────────────────────────────────────────────────────┐
│                    LEVEL SELECT SCREEN                          │
├────────────────────────────────────────────────────────────────┤
│  System Design Simulator                                       │
│                                                                 │
│  [View Tutorial] ← Can re-access anytime                       │
│  ───────────────────────────────────────────────────────────  │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │ Level 1: Local Blog                        [UNLOCKED]   │  │
│  │ Handle 10 users, Budget: $10                            │  │
│  │ Scenario: Blog for friends                              │  │
│  │                                         [Play →]        │  │
│  └─────────────────────────────────────────────────────────┘  │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │ Level 2: Growing Blog                      [LOCKED]     │  │
│  │ Complete Level 1 to unlock                              │  │
│  └─────────────────────────────────────────────────────────┘  │
│                                                                 │
│  ... more levels ...                                           │
│                                                                 │
└────────────────────────┬───────────────────────────────────────┘
                         │ Click [Play]
                         ▼
┌────────────────────────────────────────────────────────────────┐
│                        GAME SCREEN                              │
├────────────────────────────────────────────────────────────────┤
│ ┌──────────┬────────────────────────────┬─────────────────┐   │
│ │ TOOLBOX  │      CANVAS AREA           │  METRICS PANEL  │   │
│ │          │                            │                 │   │
│ │ [API]    │   ┌─────┐     ┌─────┐    │ Level: Blog     │   │
│ │ [DB]     │   │ API │────▶│ DB  │    │                 │   │
│ │ [Cache]  │   └─────┘     └─────┘    │ Objectives:     │   │
│ │ [LB]     │                           │ • Latency<500ms │   │
│ │ [CDN]    │   (Animated particles     │ • Uptime>95%    │   │
│ │          │    flow on connections)   │ • Cost<$10      │   │
│ │ ────────│                            │                 │   │
│ │ Quick    │                            │ Real-time:      │   │
│ │ Guide:   │                            │ Requests: 1234  │   │
│ │ • Click  │                            │ Latency: 45ms   │   │
│ │   to add │                            │ Cost: $2.50     │   │
│ │ • Right- │                            │                 │   │
│ │   click  │                            │ Colors:         │   │
│ │   to     │                            │ Green=Healthy   │   │
│ │   connect│                            │ Red=Overloaded  │   │
│ │          │                            │                 │   │
│ │ [? Help] │◀── Shows scenario details  │                 │   │
│ └──────────┴────────────────────────────┴─────────────────┘   │
│ ┌────────────────────────────────────────────────────────┐    │
│ │ [Start] [Stop] [Submit Solution] [Back to Levels]     │    │
│ └────────────────────────────────────────────────────────┘    │
│                                                                 │
│  Click [? Help]:                                               │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │ Scenario: You've built a blog for friends               │  │
│  │ Challenge: Handle 10 concurrent readers                 │  │
│  │                                                          │  │
│  │ Component Guide:                                         │  │
│  │ • API (Blue): Handles requests                          │  │
│  │ • DB (Purple): Stores data                              │  │
│  │                                                          │  │
│  │ Valid Connections:                                       │  │
│  │ API → Database                                           │  │
│  │                                                          │  │
│  │                                      [Close]            │  │
│  └─────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘

```

## Key Features

### 1. First Launch Experience
```
New User → Welcome Screen → 7-Page Tutorial → Level Select
   │                                               │
   └─── Learns everything before playing ─────────┘
```

### 2. Returning User Experience
```
Returning User → Level Select (directly) → Can access tutorial anytime
   │                     │
   │                     └─── [View Tutorial] button
   └─── No forced tutorial on subsequent launches
```

### 3. In-Game Help
```
Playing Level → Click [? Help] → Scenario + Component Guide → Close
   │                                                             │
   └───────────── Continues playing ────────────────────────────┘
```

## Content Breakdown

### Tutorial Pages (7 total):
1. **Welcome** - Game intro and goals
2. **Scenarios** - All 5 levels explained with context
3. **Components 1** - API Server & Database
4. **Components 2** - Cache & Load Balancer  
5. **Components 3** - CDN & Networking
6. **How to Play** - Controls and interactions
7. **Scoring** - Requirements, bonuses, tips

### Each Page Contains:
- Clear title and section headers
- Practical information
- Real examples
- Visual descriptions (colors, patterns)
- Concise, scannable text

### Navigation:
- Previous/Next buttons
- Skip Tutorial (always available)
- Page indicator (Page X of 7)
- Smart button states (disabled when appropriate)

## User Benefits

### Before Tutorial:
❌ Confused about what to do
❌ Doesn't understand scenarios
❌ Unclear on component purpose
❌ Doesn't know how to connect things
❌ Unsure about scoring

### After Tutorial:
✅ Clear understanding of game purpose
✅ Knows each level's scenario and context
✅ Understands all components deeply
✅ Can connect components correctly
✅ Knows how to optimize for scoring

## Implementation Quality

### Code Quality:
- Clean separation of concerns
- Reusable components
- State management with preferences
- Graceful error handling
- Consistent styling

### User Experience:
- Progressive disclosure
- Can skip if experienced
- Available anytime
- Non-blocking
- Context-sensitive help

### Educational Design:
- Scenario-first approach
- Why before how
- Building knowledge progressively
- Practical focus
- Clear examples

## Result

The System Design Simulator now provides a **world-class onboarding experience** that:
- Teaches system design concepts clearly
- Provides real-world context for each level
- Explains all components thoroughly
- Shows how to play effectively
- Supports both new and experienced users

This transforms the game from a puzzle to solve into a **guided learning experience**!
