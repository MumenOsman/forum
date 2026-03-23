package database

import (
	"database/sql"
	"fmt"
	"os"
)

// SeedDemoData is the main entry point for populating the database with demo content.
func SeedDemoData(db *sql.DB) {
	fmt.Println("Seeding expanded demo data...")
	seedUsers(db)
	seedPosts(db)
	seedComments(db)
}

func seedUsers(db *sql.DB) {
	users := []struct {
		id       string
		email    string
		username string
		password string
		about    string
	}{
		{
			id:       "user-elara",
			email:    "test@test.com",
			username: "Elara",
			password: "$2a$10$ejeGRPm7B7DabU/Hist3nOt5TO8A34T/pJxNBsEY56IxCJ4rQAHBG",
			about:    "The Worldbuilder. Obsessed with deep lore, complex magic systems, and the subtle mechanics of secondary worlds. Fantasy and Sci-Fi are my oxygen.",
		},
		{
			id:       "user-marcus",
			email:    "marcus@inkwell.com",
			username: "Marcus",
			password: "$2a$10$ejeGRPm7B7DabU/Hist3nOt5TO8A34T/pJxNBsEY56IxCJ4rQAHBG",
			about:    "The Historian. I value accuracy above all else. If your historical fiction gets a date wrong, we're going to have a long conversation about the Tudors.",
		},
		{
			id:       "user-chloe",
			email:    "chloe@inkwell.com",
			username: "Chloe",
			password: "$2a$10$ejeGRPm7B7DabU/Hist3nOt5TO8A34T/pJxNBsEY56IxCJ4rQAHBG",
			about:    "The Escapist. Here for the chemistry, the twists, and the happy endings! Romance and Mystery are my world. Life is too short for boring books! ✨",
		},
		{
			id:       "user-julian",
			email:    "julian@inkwell.com",
			username: "Julian",
			password: "$2a$10$ejeGRPm7B7DabU/Hist3nOt5TO8A34T/pJxNBsEY56IxCJ4rQAHBG",
			about:    "The Pragmatist. I read for efficiency. If a book doesn't offer actionable advice or a new perspective on reality, it's a DNF.",
		},
		{
			id:       "user-sarah",
			email:    "sarah@inkwell.com",
			username: "Sarah",
			password: "$2a$10$ejeGRPm7B7DabU/Hist3nOt5TO8A34T/pJxNBsEY56IxCJ4rQAHBG",
			about:    "The Aspiring Author. Analyzing pacing, structure, and character arcs as I build my own debut. Mystery and Fiction are my masterclasses.",
		},
	}

	for _, u := range users {
		_, err := db.Exec(`INSERT OR IGNORE INTO users (id, email, username, password, about_me) VALUES (?, ?, ?, ?, ?)`,
			u.id, u.email, u.username, u.password, u.about)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to seed user %s: %v\n", u.username, err)
		}
	}
}

func seedPosts(db *sql.DB) {
	posts := []struct {
		id       string
		userID   string
		title    string
		content  string
		category string
		likes    int
		dislikes int
	}{
		// Elara's Posts
		{
			id:       "post-elara-1",
			userID:   "user-elara",
			title:    "The Name of the Wind — A Reader's First Encounter",
			category: "Fantasy",
			likes:    12,
			dislikes: 0,
			content: `I picked up Patrick Rothfuss's "The Name of the Wind" after hearing endless praise for its prose, and I must say, the rumors are true. Rothfuss writes with a rare, lyrical precision that makes even a simple description of a field feel like a polished gemstone. The story centers on Kvothe, a man who has passed into legend while still alive, now hiding as a humble innkeeper and finally ready to tell his true story. 

The frame narrative—a story within a story—is a classic device, but here it feels exceptionally poignant. We see the contrast between the tired, broken man behind the bar and the brilliant, arrogant, and deeply talented boy he once was. The magic system, Sympathy, is one of the most intellectually satisfying I've ever encountered. It's grounded in physical laws of energy transfer and mental focus, making it feel less like "magic" and more like a misunderstood science. 

However, I should warn potential readers that the pacing is remarkably slow. This isn't a book of epic battles and constant action. It's a slow burn that focuses on the minutiae of Kvothe's struggles—paying for tuition, finding a place to sleep, and the agonizingly slow accumulation of knowledge. Some might find the level of detail regarding his poverty frustrating, but for me, it anchored the high-fantasy elements in a harsh, relatable reality. It’s a masterpiece of character voice, though it leaves you desperate for the next installment that may never come.`,
		},
		{
			id:       "post-elara-2",
			userID:   "user-elara",
			title:    "Dune: Masterpiece or Dry Sand?",
			category: "Sci-Fi",
			likes:    8,
			dislikes: 1,
			content: `Frank Herbert's "Dune" is often called the "Lord of the Rings" of science fiction, and after my most recent reread, I understand why. The scope of the political ecology Herbert built is staggering. He doesn't just give us a planet; he gives us a complex web of religion, economics (the Spice!), and biology that feels ancient and lived-in.

The world of Arrakis is brutal, and Herbert’s descriptions of the "stillsuits" and water-scarcity culture are immersive. You can almost feel the grit in your teeth. I particularly love the Bene Gesserit—their "magic" is actually a high-level mastery of psychology and physiology, which fits perfectly into Elara's preference for magic systems with clear internal logic.

But let's be honest: the prose can be as dry as the Arrakis desert. Herbert isn't interested in making you fall in love with his characters' personalities. Paul Atreides often feels more like a philosophical concept or a historical inevitability than a person. The dialogue is frequently formal and stilted, serving the plot and the grand themes of messianic figures and environmental collapse rather than human connection. It’s a difficult mountain to climb, and while the view from the top is spectacular, I can see why many readers turn back before the first hundred pages.`,
		},
		{
			id:       "post-elara-3",
			userID:   "user-elara",
			title:    "Advice: How to Build a Magic System that Makes Sense",
			category: "Fantasy",
			likes:    25,
			dislikes: 0,
			content: `One of the most common pitfalls I see in modern fantasy is the "Deus Ex Machina" magic system—where a character suddenly discovers a new power just in time to escape a corner. To me, a magic system only feels real when it has limitations. Here’s a quick guide on building systems that actually stick with the reader:

1. **Hard vs. Soft Magic**: Are you writing like Brandon Sanderson (Hard Magic), where the reader knows the exact rules and costs? Or like Tolkien (Soft Magic), where the power is mysterious and atmospheric? Both work, but you cannot switch mid-stream to solve a plot point. If you use Soft Magic, it should never be the primary tool for solving a conflict.

2. **Establish the Cost**: Magic should never be free. Does it drain physical stamina? Does it require expensive material components? Does it cost a piece of the user's memory? When magic has a high price, every use becomes a high-stakes decision for the character.

3. **Limitations Breed Creativity**: Instead of asking what your magic *can* do, ask what it *cannot* do. A character who can only move metal is much more interesting than a character who can simply "do magic." Limitations force characters to think their way out of problems, which makes for better storytelling.

4. **Internal Consistency**: Once you set a rule, you must never break it. If your magic depends on moonlight, don't let a character use it in a cave unless you've spent the whole book explaining a way to store lunar energy. Consistency builds trust with your audience.`,
		},
		{
			id:       "post-elara-4",
			userID:   "user-elara",
			title:    "Series Review: The Expanse (Books 1-3)",
			category: "Sci-Fi",
			likes:    15,
			dislikes: 2,
			content: `I’ve just finished the first trilogy of "The Expanse" (Leviathan Wakes, Caliban's War, and Abaddon's Gate), and I’m completely hooked. What sets this series apart from other space operas is the commitment to "hard" science—or at least the illusion of it. Gravity isn't just a setting; it's a plot point. The effects of high-G burns on the human body and the biological differences between "Belters" and "Inners" create a fascinating sociological layer.

The political tension between Earth, Mars, and the Belt feels incredibly grounded. It’s not a simple "good vs. evil" conflict. Every faction has a legitimate grievance and a logical goal. The Rocinante crew—Holden, Naomi, Amos, and Alex—feel like a real family, warts and all. Their evolution from accidental outlaws to a cohesive political force is one of the best character-driven arcs in modern Sci-Fi.

If you like your science fiction with a heavy dose of realism, political maneuvering, and a touch of cosmic horror (the protomolecule is truly unsettling), start this series immediately. It’s the closest thing we have to a modern-day masterpiece in the genre.`,
		},

		// Marcus's Posts
		{
			id:       "post-marcus-1",
			userID:   "user-marcus",
			title:    "Wolf Hall: Bringing Thomas Cromwell to Life",
			category: "Historical Fiction",
			likes:    9,
			dislikes: 0,
			content: `Hilary Mantel's "Wolf Hall" is, in my professional opinion, the pinnacle of modern historical fiction. Mantel achieves something nearly impossible: she takes one of the most reviled figures in English history—Thomas Cromwell—and makes him deeply, magnetically human. 

The most striking stylistic choice is the use of the present tense. It strips away the "dustiness" of history and makes the political maneuvering of Henry VIII’s court feel urgent and dangerous. We are seeing the 1520s not as a series of facts in a textbook, but as a lived experience. Cromwell is portrayed as a man of immense intellect and practicality, a "man of business" navigating a world of fickle kings and rigid social hierarchies.

The historical research is meticulous. Mantel doesn't just get the dates and names right; she understands the worldview of the time. The religious anxieties of the Reformation aren't just background noise—they are the driving force of the characters' lives. It’s a dense read, and you need a basic grasp of the period to fully appreciate it, but for anyone who values historical accuracy married to incredible prose, this is the gold standard.`,
		},
		{
			id:       "post-marcus-2",
			userID:   "user-marcus",
			title:    "Chernow's Hamilton: More Than Just a Musical",
			category: "Biography",
			likes:    14,
			dislikes: 1,
			content: `While many are introduced to Alexander Hamilton through the stage show, Ron Chernow’s 800-page biography is where the real complexity of the man lies. Chernow doesn't shy away from Hamilton's flaws—his volatility, his arrogance, and his almost pathological need to respond to every criticism in writing (the Reynolds Pamphlet being the prime example of his self-destruction).

What I found most compelling was the detailed account of his childhood in the West Indies. To understand Hamilton’s drive, you must understand the grinding poverty and trauma he escaped. Chernow argues that this background gave him a unique perspective among the Founding Fathers, allowing him to see the necessity of a modern financial system and a strong central government in a way that the aristocratic Jefferson never could.

The book is an exhaustive study of the birth of the American administrative state. It can be a slog through the sections on the Coast Guard and the central bank, but Chernow makes the intellectual battles of the time feel as personal as a duel. It is a portrait of a man who was his own greatest asset and his own worst enemy.`,
		},
		{
			id:       "post-marcus-3",
			userID:   "user-marcus",
			title:    "Advice: Balancing Fact and Fiction in Historical Settings",
			category: "Historical Fiction",
			likes:    18,
			dislikes: 0,
			content: `Writing historical fiction is a balancing act that many fail. Lean too far into "Fact" and you have a dry academic paper; lean too far into "Fiction" and you risk alienating the very readers attracted to the era. Here is my advice for finding the middle ground:

1. **The 'In-Between' Spaces**: The best place for fiction is in the gaps left by the historical record. We know *what* happened in 1536, but we don't always know *what was said* behind closed doors. That is where your story lives.

2. **Research for Atmosphere, Not Just Facts**: Don't just look up dates. Research what people ate, how their clothes felt, what the air smelled like in a crowded 18th-century street. These sensory details anchor the reader more than a list of battle statistics.

3. **The "Enough" Rule**: Stop researching when you start finding the same information over and over. You don't need to be an expert on medieval tax law to write a story about a squire, but you do need to know enough to not make him sound like a 21st-century teenager.

4. **Honesty About Inaccuracy**: If you must change a date for the sake of the plot, own it in an Author's Note. Most historical readers will forgive a minor shift if the story is compelling, but they will never forgive being lied to.`,
		},
		{
			id:       "post-marcus-4",
			userID:   "user-marcus",
			title:    "Isaacson's Steve Jobs: The Cost of Genius",
			category: "Biography",
			likes:    5,
			dislikes: 3,
			content: `Walter Isaacson’s biography of Steve Jobs is a fascinating, if deeply uncomfortable, read. Isaacson portrays Jobs as a visionary who changed the world, but he also documents a level of interpersonal cruelty that is hard to stomach. From his denial of paternity of his first child to his "reality distortion field" used to berate employees, the book raises a critical question: was the toxicity a necessary byproduct of the innovation?

Personally, I'm skeptical. While his drive for perfection led to the iPhone, his refusal to seek traditional medical treatment for his cancer—relying instead on "alternative" methods until it was too late—shows the dark side of that same stubbornness. 

Isaacson’s writing is straightforward and journalistic, which is appropriate for a man who famously insisted on "no secrets" for his biography. It’s a well-structured book that traces the rise, fall, and resurrection of Apple, but by the end, I felt less admiration for the man and more pity for those who had to live and work in his orbit. A high cost for genius, indeed.`,
		},

		// Chloe's Posts
		{
			id:       "post-chloe-1",
			userID:   "user-chloe",
			title:    "Why Pride and Prejudice is Still the Blueprint",
			category: "Romance",
			likes:    30,
			dislikes: 1,
			content: `Oh my gosh, you guys, I just reread Pride and Prejudice (for like the 50th time!) and I am still absolutely swooning over Mr. Darcy! 😍 

Jane Austen was truly the queen of the "enemies-to-lovers" trope before it even had a name. The chemistry between Lizzy and Darcy isn't about physical touch—it's all in those longing looks, the sharp banter, and that disastrous first proposal! The way they both have to grow and admit they were wrong (Lizzy about her prejudice and Darcy about his pride) is just peak character development. 

Austen's wit is so sharp that some of the insults still sting 200 years later. "Tolerable, but not handsome enough to tempt me"—ouch! But when he finally says, "You must allow me to tell you how ardently I admire and love you," my heart just melts every single time! It’s the perfect blend of societal drama, family chaos, and a love story that feels earned. If you haven't read it because you think it's "stuffy," please give it a try! It’s basically the original rom-com! ✨☕️`,
		},
		{
			id:       "post-chloe-2",
			userID:   "user-chloe",
			title:    "Gone Girl: The Twist That Broke My Brain",
			category: "Mystery",
			likes:    22,
			dislikes: 4,
			content: `Okay, let’s talk about Gone Girl by Gillian Flynn. 😱 I’m not going to post any spoilers here because everyone deserves to experience this for themselves, but oh my god!! 

The unreliability of the narrators is just... wow. One minute you're feeling so sorry for Nick, and the next you're reading Amy's diary and you're not sure who to believe. And then the middle of the book happens and everything you thought you knew gets flipped upside down! 🔄

It’s a very dark look at a toxic marriage and how well we actually know the people we love. I stayed up until 3 AM finishing this because I literally could not breathe until I knew how it ended. It’s definitely not a "happy" book, and some parts are really disturbing, but as a mystery, it is absolutely flawless. The pacing is intense and the writing is like a knife. Just a warning: you might want to hug someone you trust after finishing this one! 😅🕵️‍♀️`,
		},
		{
			id:       "post-chloe-3",
			userID:   "user-chloe",
			title:    "Advice: Structuring the Perfect 'Slow Burn'",
			category: "Romance",
			likes:    19,
			dislikes: 0,
			content: `Is there anything better in a book than a perfectly executed slow burn? That feeling of "will they, won't they" that keeps you up all night turning pages? 💖 Here are my tips for making the tension absolutely unbearable (in a good way!):

1. **Forced Proximity**: Put them in a situation where they can’t escape each other! A long road trip, stuck in an elevator, or (my favorite) the "only one bed" trope! It forces them to interact when they aren't ready for it. 🏨

2. **Micro-Tensions**: It’s not about the big kiss! It’s about the "accidental" hand brush, the look from across the room, or the way he remembers her favorite coffee order. Those tiny details build the foundation for the big moments.

3. **External Stakes**: Give them a reason why they *can't* be together right now. Is it a family feud? A career rivalry? A secret they're hiding? The external conflict keeps them apart while the internal feelings are pulling them together.

4. **Delay the Payoff**: Don't give in too early! Let the reader get just a little bit frustrated. If they kiss on page 50, the tension is gone. Wait until they've almost lost each other before the big emotional payoff. Trust me, it makes the ending so much sweeter! 🍭✨`,
		},
		{
			id:       "post-chloe-4",
			userID:   "user-chloe",
			title:    "The Maidens by Alex Michaelides – A Letdown?",
			category: "Mystery",
			likes:    7,
			dislikes: 2,
			content: `I loved "The Silent Patient" so much that I was counting down the days until "The Maidens" came out, but honestly? I’m so disappointed. 😔

The atmospheric setting of Cambridge University was great—it felt very dark academia—but the plot just didn't hold up for me. The main character, Mariana, made so many questionable decisions that it was hard to stay on her side. And the "big twist" at the end? It felt completely unearned. Like, there was almost no foreshadowing, so instead of feeling like a "Gotcha!" moment, it just felt like the author made it up at the last second to be shocking. 

Compared to the brilliant psychological layering of his first book, this one felt a bit messy and rushed. I still love Michaelides' writing style, but the mystery itself didn't have that "click" moment where everything suddenly makes sense. Am I the only one who felt this way?! Let me know! 📚💔`,
		},

		// Julian's Posts
		{
			id:       "post-julian-1",
			userID:   "user-julian",
			title:    "Thinking, Fast and Slow: Changing How I Decide",
			category: "Non-Fiction",
			likes:    11,
			dislikes: 1,
			content: `Daniel Kahneman’s "Thinking, Fast and Slow" is required reading for anyone interested in why they make mistakes. The core concept is simple: we have two systems of thought. 

**System 1** is fast, instinctive, and emotional. It’s what lets us recognize a face or drive a familiar route without thinking. **System 2** is slower, more deliberative, and logical. It handles complex math and difficult decisions.

The problem, as Kahneman explains, is that System 1 is lazy and likes to take shortcuts, which leads to cognitive biases like the "availability heuristic" or "loss aversion." Since reading this, I’ve started pausing before making any significant purchase or business decision to ask: "Is this System 1's impulse or System 2's analysis?" 

It’s a dense, academic book, and you can skip the statistics sections if you're in a rush, but the takeaways on human irrationality are invaluable. It has fundamentally changed how I view my own brain.`,
		},
		{
			id:       "post-julian-2",
			userID:   "user-julian",
			title:    "Advice: How to Read 50 Books a Year Without Burnout",
			category: "General",
			likes:    45,
			dislikes: 2,
			content: `People often ask me how I manage to read so much while running a business. It’s not about "speed reading" (which is mostly a myth); it’s about systems. Here is my 3-step process:

1. **The DNF Rule**: Life is too short for bad books. If you aren't hooked by page 50, stop reading. Do not let a boring book create a "reading bottleneck" where you don't read anything for a month because you feel guilty about not finishing it. Move on.

2. **Audiobooks for "Dead Time"**: I "read" while driving, at the gym, and doing chores. Non-fiction is particularly suited to audio. If you spend 30 minutes a day commuting, that’s 2.5 hours a week of reading time you didn't have before.

3. **The Morning 20**: Read for 20 minutes before you check your phone. No emails, no news, no social media. Just you and a book. It sets a productive tone for the day and ensures that even on your busiest days, you’ve made progress.

Consistency beats intensity every time. You don't need a 4-hour block of time; you just need to stop wasting the 10-minute blocks you already have.`,
		},
		{
			id:       "post-julian-3",
			userID:   "user-julian",
			title:    "Atomic Habits: Is it overhyped?",
			category: "Non-Fiction",
			likes:    33,
			dislikes: 5,
			content: `James Clear’s "Atomic Habits" is the best-selling habit book for a reason: it’s incredibly actionable. His "Four Laws of Behavior Change" (Make it Obvious, Make it Attractive, Make it Easy, Make it Satisfying) are solid psychological principles packaged for the layman.

However, I have a slight critique. The book could easily have been a 20-page long-form essay. There is a lot of repetition and fluff used to stretch the concept into a 300-page hardcover. Once you understand the concept of "identity-based habits" and "habit stacking," you’ve basically got the meat of the book. 

Is it worth reading? Yes, because the repetition actually helps the concepts stick. But if you're short on time, you can get 90% of the value by reading a detailed summary and applying the laws to one small habit tomorrow. Good, but efficiency-wise, it’s a bit bloated.`,
		},
		{
			id:       "post-julian-4",
			userID:   "user-julian",
			title:    "My Top 5 Books of the Decade",
			category: "General",
			likes:    20,
			dislikes: 0,
			content: `As we cross the halfway mark of the 2020s, I wanted to share the five books from the last decade that have had the biggest impact on my worldview. Most of these are non-fiction, but they all share a commitment to objective truth and actionable insights.

1. **The Selfish Gene (Updated Edition)** - Dawkins. The foundation of understanding biology and human competition.
2. **Sapiens** - Yuval Noah Harari. A brilliant, if sometimes controversial, look at the stories that hold our society together.
3. **The Big Short** - Michael Lewis. The best explanation of financial systems and the danger of groupthink.
4. **Bad Blood** - John Carreyrou. A masterclass in investigative journalism and the dangers of the Silicon Valley "fake it 'til you make it" culture.
5. **Flow** - Mihaly Csikszentmihalyi. The science of optimal experience. This should be required reading for every professional.

What are yours? Keep it brief—I value the quality of the list over the quantity.`,
		},

		// Sarah's Posts
		{
			id:       "post-sarah-1",
			userID:   "user-sarah",
			title:    "Advice: How to Write a Novel - Outlining vs. Pantsing",
			category: "Fiction",
			likes:    35,
			dislikes: 1,
			content: `The eternal debate among writers: are you an **Outliner** (Plotter) or a **Pantser** (Discovery Writer)? As someone currently in the trenches of my first draft, I’ve realized that most of us fall somewhere in the middle. 

**Outlining** gives you a roadmap. You know your beats, your twists, and your ending before you type "Chapter 1." It prevents that dreaded "middle-of-the-book slump" because you always know what needs to happen next. But some say it kills the magic of discovery.

**Pantsing** (writing by the seat of your pants) allows characters to surprise you. The story feels organic because even the author doesn't know what’s coming. The downside? You often end up with 100,000 words that go nowhere and require a massive, painful rewrite.

My advice? **Plantser** it. Outline the major "tentpole" moments (The Inciting Incident, The Midpoint, The Climax), but leave the path between them open to exploration. It gives you the safety of a structure with the freedom to let the characters breathe. What’s your style? I’m finding that the older I get, the more I crave the safety of a plan!`,
		},
		{
			id:       "post-sarah-2",
			userID:   "user-sarah",
			title:    "The Secret History: Beautiful and Toxic",
			category: "Fiction",
			likes:    16,
			dislikes: 0,
			content: `Donna Tartt's "The Secret History" is essentially the reason the "Dark Academia" aesthetic exists. Her prose is lush, atmospheric, and incredibly immersive. You feel like you’re right there in that snowy Vermont college, drinking tea and reading Greek classics with the most pretentious, fascinating people you've ever met.

But what I find most interesting as a writer is the ensemble cast. None of them are "good" people. They are elitist, detached, and eventually, murderous. Yet, Tartt makes you care about their descent. She uses a "reverse mystery" structure—we know they killed Bunny on page one, so the tension isn't "who did it," but "why they did it" and "how they fall apart afterward."

It’s a masterclass in building tension through character dynamics rather than external action. The way she uses classical Greek themes of fate and hubris to mirror the modern-day actions of the students is brilliant. If you can stomach a group of characters who are deeply flawed and often unlikeable, the prose alone is worth the price of admission.`,
		},
		{
			id:       "post-sarah-3",
			userID:   "user-sarah",
			title:    "Advice: Dropping Red Herrings Like a Pro",
			category: "Mystery",
			likes:    28,
			dislikes: 0,
			content: `In a good mystery, the reader wants to be fooled, but they also want to feel like they *could* have solved it if they were just a little smarter. Enter the **Red Herring**. Here’s how to use them without making your reader feel cheated:

1. **The 'Loud' Clue**: Make a red herring look like the most important thing in the room. If a character is suspiciously cleaning a knife when the detective walks in, the reader will focus on it—meanwhile, the real clue is something small and unnoticed in the background.

2. **Justify the Suspicion**: A red herring shouldn't just be a random lie. If a character is acting weirdly, they should have a *reason* for it that isn't the murder. Maybe they’re hiding an affair, or a gambling debt. This makes them a valid suspect for the reader while keeping the plot grounded.

3. **Follow the 'Fair Play' Rule**: You must give the reader all the information they need to spot the red herring. If the knife is the red herring, there should be a throwaway line later about the character preparing for a dinner party. 

4. **Synchronicity**: Drop the real clue and the red herring at the same time. The reader's brain will usually latch onto the more dramatic of the two, letting the real truth hide in plain sight. Happy plotting! 🖋️🕵️‍♂️`,
		},
		{
			id:       "post-sarah-4",
			userID:   "user-sarah",
			title:    "Tomorrow, and Tomorrow, and Tomorrow — A Masterclass in Character",
			category: "Fiction",
			likes:    21,
			dislikes: 2,
			content: `I’ve just finished Gabrielle Zevin’s "Tomorrow, and Tomorrow, and Tomorrow," and my heart is essentially in pieces. 😭 This isn't just a book about video games; it's a book about the creative process and the agonizing, beautiful complexity of platonic love.

Following Sam and Sadie (and Marx, the unsung hero!) over thirty years was like watching a real life unfold. Zevin captures the way friendships drift, collide, and reshape themselves with such painful accuracy. As a writer, I was floored by how she depicted the "creative partnership"—that intense, often toxic, always inspiring connection that happens when two people build a world together.

Even if you aren't a "gamer," the sections describing the games they build are so evocative. They mirror the characters' internal states perfectly. It’s a masterclass in how to show, not tell. It’s a long book, and it doesn't always go where you want it to, but it feels *true*. If you want a deep dive into what it means to love someone you can't quite be with and haven't quite lost, this is it. Brilliant.`,
		},
	}

	for _, p := range posts {
		// Insert Post
		_, err := db.Exec(`INSERT OR IGNORE INTO posts (id, user_id, title, content, likes, dislikes) VALUES (?, ?, ?, ?, ?, ?)`,
			p.id, p.userID, p.title, p.content, p.likes, p.dislikes)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to seed post %s: %v\n", p.title, err)
			continue
		}

		// Link to Category
		var catID int
		err = db.QueryRow(`SELECT id FROM categories WHERE name = ?`, p.category).Scan(&catID)
		if err == nil {
			_, _ = db.Exec(`INSERT OR IGNORE INTO post_categories (post_id, category_id) VALUES (?, ?)`,
				p.id, catID)
		}
	}
}

func seedComments(db *sql.DB) {
	comments := []struct {
		id      string
		postID  string
		userID  string
		content string
		likes   int
	}{
		// Comments on Post 1 (The Name of the Wind)
		{id: "c-1", postID: "post-elara-1", userID: "user-marcus", content: "I quite enjoyed the 'interlude' chapters at the Waystone Inn. The framing of an old legend hiding as a simple tavern-keeper is a classic historical device, handled with great dignity here.", likes: 3},
		{id: "c-2", postID: "post-elara-1", userID: "user-sarah", content: "That frame-narrative is such a masterclass in tension! We know exactly where he ends up, but the 'how' is what keeps the pages turning. Rothfuss is brilliant at withholding just enough.", likes: 5},
		
		// Comments on Post 2 (Dune)
		{id: "c-3", postID: "post-elara-2", userID: "user-julian", content: "I agree, the prose was a slog, but the lessons on long-term ecological planning and the danger of charismatic leaders are still incredibly relevant for anyone in leadership today.", likes: 4},
		
		// Comments on Post 3 (Magic Systems)
		{id: "c-4", postID: "post-elara-3", userID: "user-sarah", content: "Setting costs is something I've been struggling with. Thanks for the tip—I think making my protagonist lose a literal memory for each spell would add that stakes I was missing!", likes: 8},
		{id: "c-5", postID: "post-elara-3", userID: "user-chloe", content: "This is so helpful! Can I use a 'cost' like this in a paranormal romance, or is that too grim? I want the love to be the focus, but the magic should feel real too! ✨", likes: 2},
		
		// Comments on Post 5 (Wolf Hall)
		{id: "c-6", postID: "post-marcus-1", userID: "user-elara", content: "The way Cromwell manipulates the court feels so much like the political intrigue in some of my favorite high-fantasy series. Do you think George R.R. Martin was inspired by the Tudors?", likes: 6},
		
		// Comments on Post 6 (Hamilton)
		{id: "c-7", postID: "post-marcus-2", userID: "user-julian", content: "Hamilton's work ethic was truly insane. Chernow makes a great point—he was a master of high-output management before the term even existed. It’s basically a productivity manual.", likes: 4},
		
		// Comments on Post 7 (Historical Accuracy)
		{id: "c-8", postID: "post-marcus-3", userID: "user-sarah", content: "I'm definitely taking notes on the 'In-Between' spaces strategy. It feels much less intimidating than trying to rewrite literal history!", likes: 5},
		
		// Comments on Post 8 (Steve Jobs)
		{id: "c-9", postID: "post-marcus-4", userID: "user-julian", content: "While I don't condone the behavior, it's hard to argue with the results. Jobs understood the intersection of technology and the liberal arts better than anyone in history.", likes: 7},
		{id: "c-10", postID: "post-marcus-4", userID: "user-chloe", content: "Honestly, I couldn't even finish the book because he was just too mean to everyone! 😭 I don't care how cool the iPhone is, you have to be a nice person!", likes: 12},
		
		// Comments on Post 9 (Pride and Prejudice)
		{id: "c-11", postID: "post-chloe-1", userID: "user-marcus", content: "The historical societal constraints were actually quite realistic for the Regency era. The desperation for a 'good match' was a matter of survival, not just romance.", likes: 9},
		{id: "c-12", postID: "post-chloe-1", userID: "user-sarah", content: "As a writer, I'm always analyzing how Austen manages the 'prejudice' shift. It's so subtle! You don't realize Lizzy is changing her mind until she already has.", likes: 11},
		
		// Comments on Post 10 (Gone Girl)
		{id: "c-13", postID: "post-chloe-2", userID: "user-sarah", content: "The structural shift in the middle is absolutely legendary. I still remember the physical chill I got when I turned that page. Truly a masterclass in tension.", likes: 14},
		{id: "c-14", postID: "post-chloe-2", userID: "user-elara", content: "It's brilliantly written, but honestly? It was a bit too dark for my taste. I prefer my villains to have a bit more... magic, perhaps? Or at least some honor.", likes: 3},
		
		// Comments on Post 11 (Slow Burn)
		{id: "c-15", postID: "post-chloe-3", userID: "user-sarah", content: "How do you keep the middle from sagging when nothing 'big' is happening yet? I always feel like I'm just treading water until the confession!", likes: 6},
		
		// Comments on Post 13 (Thinking Fast and Slow)
		{id: "c-16", postID: "post-julian-1", userID: "user-marcus", content: "I can see so many historical blunders that were clearly System 1 errors. Napoleon's march on Moscow comes to mind immediately. Excessive overconfidence!", likes: 5},
		
		// Comments on Post 14 (50 Books a Year)
		{id: "c-17", postID: "post-julian-2", userID: "user-elara", content: "I actually disagree with the DNF rule! Sometimes a book is a struggle for a reason, and the reward is in the completion. I've never abandoned a book and I never will!", likes: 4},
		{id: "c-18", postID: "post-julian-2", userID: "user-chloe", content: "I love the audiobook tip! I listen to all my thrillers while I'm doing my skincare routine and it's the best part of my day! 🌸✨", likes: 8},
		{id: "c-19", postID: "post-julian-2", userID: "user-sarah", content: "I'm officially adopting the DNF rule starting today. I have three books on my desk that I've been avoiding for months—time to let them go!", likes: 10},
		
		// Comments on Post 15 (Atomic Habits)
		{id: "c-20", postID: "post-julian-3", userID: "user-marcus", content: "I actually disagree that it's bloated. From an educational standpoint, the repetition is vital for internalizing the laws of behavior. You have to hear it multiple times for it to stick.", likes: 6},
		
		// Comments on Post 16 (Top 5 Books)
		{id: "c-21", postID: "post-julian-4", userID: "user-chloe", content: "Julian, please!! This list is too serious! 😭 I demand that you read at least one fiction book this year and report back here!", likes: 15},
		
		// Comments on Post 17 (Outlining vs Pantsing)
		{id: "c-22", postID: "post-sarah-1", userID: "user-elara", content: "For me, the worldbuilding alone requires heavy outlining. If I don't know the geography and the magic rules beforehand, the story just collapses into a mess.", likes: 7},
		{id: "c-23", postID: "post-sarah-1", userID: "user-chloe", content: "I'm a total pantser! Romance needs that organic chemistry, and I feel like if I plan it too much, it just feels fake. Let the hearts decide! 💖", likes: 5},
		
		// Comments on Post 18 (The Secret History)
		{id: "c-24", postID: "post-sarah-2", userID: "user-marcus", content: "I quite enjoyed the classical Greek references. Tartt clearly did her homework on the Bacchae—it adds a layer of ancient tragedy to the modern setting.", likes: 8},
		
		// Comments on Post 19 (Red Herrings)
		{id: "c-25", postID: "post-sarah-3", userID: "user-chloe", content: "Yes!! This is exactly why I love thrillers so much! That moment when you realize you were looking at the wrong thing the whole time is the best! 🕵️‍♀️💖", likes: 12},
		{id: "c-26", postID: "post-sarah-3", userID: "user-elara", content: "It's quite similar to foreshadowing in fantasy. If you introduce a sword on page 10, it better be the one that kills the dragon on page 50０. Great tips!", likes: 6},
		
		// Comments on Post 20 (Tomorrow x3)
		{id: "c-27", postID: "post-sarah-4", userID: "user-julian", content: "I appreciated the realistic depiction of startup culture and the creative partnership dynamics. It's rare to see the exhaustion of building something new portrayed so accurately.", likes: 9},
	}

	for _, c := range comments {
		_, err := db.Exec(`INSERT OR IGNORE INTO comments (id, post_id, user_id, content, likes) VALUES (?, ?, ?, ?, ?)`,
			c.id, c.postID, c.userID, c.content, c.likes)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to seed comment %s: %v\n", c.id, err)
		}
	}
}
