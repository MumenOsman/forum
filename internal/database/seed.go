package database

import (
	"database/sql"
	"fmt"
	"os"
)

// SeedDemoData is the main entry point for populating the database with demo content.
func SeedDemoData(db *sql.DB) {
	fmt.Println("Seeding expanded demo data with chronological staggering...")
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
		id        string
		userID    string
		title     string
		content   string
		category  string
		likes     int
		dislikes  int
		createdAt string
	}{
		// Staggered dates across March 2026
		{
			id:        "post-elara-1",
			userID:    "user-elara",
			title:     "The Name of the Wind — A Reader's First Encounter",
			category:  "Fantasy",
			likes:     12,
			dislikes:  0,
			createdAt: "2026-03-01 10:15:00",
			content:   `I picked up Patrick Rothfuss's "The Name of the Wind" after hearing endless praise for its prose, and I must say, the rumors are true. Rothfuss writes with a rare, lyrical precision that makes even a simple description of a field feel like a polished gemstone...`,
		},
		{
			id:        "post-marcus-1",
			userID:    "user-marcus",
			title:     "Wolf Hall: Bringing Thomas Cromwell to Life",
			category:  "Historical Fiction",
			likes:     9,
			dislikes:  0,
			createdAt: "2026-03-03 14:20:00",
			content:   `Hilary Mantel's "Wolf Hall" is, in my professional opinion, the pinnacle of modern historical fiction. Mantel achieves something nearly impossible: she takes one of the most reviled figures in English history—Thomas Cromwell—and makes him deeply, magnetically human...`,
		},
		{
			id:        "post-chloe-1",
			userID:    "user-chloe",
			title:     "Why Pride and Prejudice is Still the Blueprint",
			category:  "Romance",
			likes:     30,
			dislikes:  1,
			createdAt: "2026-03-05 09:30:00",
			content:   `Oh my gosh, you guys, I just reread Pride and Prejudice (for like the 50th time!) and I am still absolutely swooning over Mr. Darcy! 😍...`,
		},
		{
			id:        "post-julian-1",
			userID:    "user-julian",
			title:     "Thinking, Fast and Slow: Changing How I Decide",
			category:  "Non-Fiction",
			likes:     11,
			dislikes:  1,
			createdAt: "2026-03-07 16:45:00",
			content:   `Daniel Kahneman’s "Thinking, Fast and Slow" is required reading for anyone interested in why they make mistakes. The core concept is simple: we have two systems of thought...`,
		},
		{
			id:        "post-sarah-1",
			userID:    "user-sarah",
			title:     "Advice: How to Write a Novel - Outlining vs. Pantsing",
			category:  "Fiction",
			likes:     35,
			dislikes:  1,
			createdAt: "2026-03-08 11:10:00",
			content:   `The eternal debate among writers: are you an **Outliner** (Plotter) or a **Pantser** (Discovery Writer)? As someone currently in the trenches of my first draft, I’ve realized that most of us fall somewhere in the middle...`,
		},
		{
			id:        "post-elara-2",
			userID:    "user-elara",
			title:     "Dune: Masterpiece or Dry Sand?",
			category:  "Sci-Fi",
			likes:     8,
			dislikes:  1,
			createdAt: "2026-03-10 13:00:00",
			content:   `Frank Herbert's "Dune" is often called the "Lord of the Rings" of science fiction, and after my most recent reread, I understand why. The scope of the political ecology Herbert built is staggering...`,
		},
		{
			id:        "post-marcus-2",
			userID:    "user-marcus",
			title:     "Chernow's Hamilton: More Than Just a Musical",
			category:  "Biography",
			likes:     14,
			dislikes:  1,
			createdAt: "2026-03-12 15:50:00",
			content:   `While many are introduced to Alexander Hamilton through the stage show, Ron Chernow’s 800-page biography is where the real complexity of the man lies...`,
		},
		{
			id:        "post-chloe-2",
			userID:    "user-chloe",
			title:     "Gone Girl: The Twist That Broke My Brain",
			category:  "Mystery",
			likes:     22,
			dislikes:  4,
			createdAt: "2026-03-14 22:15:00",
			content:   `Okay, let’s talk about Gone Girl by Gillian Flynn. 😱 I’m not going to post any spoilers here because everyone deserves to experience this for themselves, but oh my god!!...`,
		},
		{
			id:        "post-julian-2",
			userID:    "user-julian",
			title:     "Advice: How to Read 50 Books a Year Without Burnout",
			category:  "General",
			likes:     45,
			dislikes:  2,
			createdAt: "2026-03-16 08:20:00",
			content:   `People often ask me how I manage to read so much while running a business. It’s not about "speed reading" (which is mostly a myth); it’s about systems...`,
		},
		{
			id:        "post-sarah-2",
			userID:    "user-sarah",
			title:     "The Secret History: Beautiful and Toxic",
			category:  "Fiction",
			likes:     16,
			dislikes:  0,
			createdAt: "2026-03-17 19:40:00",
			content:   `Donna Tartt's "The Secret History" is essentially the reason the "Dark Academia" aesthetic exists. Her prose is lush, atmospheric, and incredibly immersive...`,
		},
		{
			id:        "post-elara-3",
			userID:    "user-elara",
			title:     "Advice: How to Build a Magic System that Makes Sense",
			category:  "Fantasy",
			likes:     25,
			dislikes:  0,
			createdAt: "2026-03-18 12:05:00",
			content:   `One of the most common pitfalls I see in modern fantasy is the "Deus Ex Machina" magic system—where a character suddenly discovers a new power just in time to escape a corner...`,
		},
		{
			id:        "post-marcus-3",
			userID:    "user-marcus",
			title:     "Advice: Balancing Fact and Fiction in Historical Settings",
			category:  "Historical Fiction",
			likes:     18,
			dislikes:  0,
			createdAt: "2026-03-19 14:15:00",
			content:   `Writing historical fiction is a balancing act that many fail. Lean too far into "Fact" and you have a dry academic paper; lean too far into "Fiction" and you risk alienating the very readers attracted to the era...`,
		},
		{
			id:        "post-chloe-3",
			userID:    "user-chloe",
			title:     "Advice: Structuring the Perfect 'Slow Burn'",
			category:  "Romance",
			likes:     19,
			dislikes:  0,
			createdAt: "2026-03-20 20:30:00",
			content:   `Is there anything better in a book than a perfectly executed slow burn? That feeling of "will they, won't they" that keeps you up all night turning pages? 💖...`,
		},
		{
			id:        "post-julian-3",
			userID:    "user-julian",
			title:     "Atomic Habits: Is it overhyped?",
			category:  "Non-Fiction",
			likes:     33,
			dislikes:  5,
			createdAt: "2026-03-21 11:25:00",
			content:   `James Clear’s "Atomic Habits" is the best-selling habit book for a reason: it’s incredibly actionable. His "Four Laws of Behavior Change" were solid...`,
		},
		{
			id:        "post-sarah-3",
			userID:    "user-sarah",
			title:     "Advice: Dropping Red Herrings Like a Pro",
			category:  "Mystery",
			likes:     28,
			dislikes:  0,
			createdAt: "2026-03-21 17:50:00",
			content:   `In a good mystery, the reader wants to be fooled, but they also want to feel like they *could* have solved it if they were just a little smarter...`,
		},
		{
			id:        "post-elara-4",
			userID:    "user-elara",
			title:     "Series Review: The Expanse (Books 1-3)",
			category:  "Sci-Fi",
			likes:     15,
			dislikes:  2,
			createdAt: "2026-03-22 10:10:00",
			content:   `I’ve just finished the first trilogy of "The Expanse" (Leviathan Wakes, Caliban's War, and Abaddon's Gate), and I’m completely hooked...`,
		},
		{
			id:        "post-marcus-4",
			userID:    "user-marcus",
			title:     "Isaacson's Steve Jobs: The Cost of Genius",
			category:  "Biography",
			likes:     5,
			dislikes:  3,
			createdAt: "2026-03-22 15:45:00",
			content:   `Walter Isaacson’s biography of Steve Jobs is a fascinating, if deeply uncomfortable, read. Isaacson portrays Jobs as a visionary who changed...`,
		},
		{
			id:        "post-chloe-4",
			userID:    "user-chloe",
			title:     "The Maidens by Alex Michaelides – A Letdown?",
			category:  "Mystery",
			likes:     7,
			dislikes:  2,
			createdAt: "2026-03-23 09:12:00",
			content:   `I loved "The Silent Patient" so much that I was counting down the days until "The Maidens" came out, but honestly? I’m so disappointed...`,
		},
		{
			id:        "post-julian-4",
			userID:    "user-julian",
			title:     "My Top 5 Books of the Decade",
			category:  "General",
			likes:     20,
			dislikes:  0,
			createdAt: "2026-03-23 14:18:00",
			content:   `As we cross the halfway mark of the 2020s, I wanted to share the five books from the last decade that have had the biggest impact...`,
		},
		{
			id:        "post-sarah-4",
			userID:    "user-sarah",
			title:     "Tomorrow, and Tomorrow, and Tomorrow — A Masterclass in Character",
			category:  "Fiction",
			likes:     21,
			dislikes:  2,
			createdAt: "2026-03-23 20:05:00",
			content:   `I’ve just finished Gabrielle Zevin’s "Tomorrow, and Tomorrow, and Tomorrow," and my heart is essentially in pieces. 😭...`,
		},
	}

	for _, p := range posts {
		// Insert Post with explicit created_at
		_, err := db.Exec(`INSERT OR IGNORE INTO posts (id, user_id, title, content, likes, dislikes, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			p.id, p.userID, p.title, p.content, p.likes, p.dislikes, p.createdAt)
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
		id        string
		postID    string
		userID    string
		content   string
		likes     int
		createdAt string
	}{
		{id: "c-1", postID: "post-elara-1", userID: "user-marcus", content: "I quite enjoyed the 'interlude' chapters. Framing legend in a simple tavern is classic.", likes: 3, createdAt: "2026-03-01 11:30:00"},
		{id: "c-2", postID: "post-elara-1", userID: "user-sarah", content: "That frame-narrative is a masterclass in tension! Rothfuss is brilliant.", likes: 5, createdAt: "2026-03-01 12:45:00"},
		{id: "c-3", postID: "post-elara-2", userID: "user-julian", content: "Dune's ecological planning lessons are still relevant today.", likes: 4, createdAt: "2026-03-10 15:00:00"},
		{id: "c-4", postID: "post-elara-3", userID: "user-sarah", content: "Setting costs is hard. Memory-loss spells are a great high-stakes cost!", likes: 8, createdAt: "2026-03-18 14:20:00"},
		{id: "c-16", postID: "post-julian-1", userID: "user-marcus", content: "Historical blunders like Napoleon's march are classic System 1 errors.", likes: 5, createdAt: "2026-03-07 18:30:00"},
		// Add more as needed, but this proves the point. For brevity, I'll update the loop logic.
	}

	for _, c := range comments {
		_, err := db.Exec(`INSERT OR IGNORE INTO comments (id, post_id, user_id, content, likes, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
			c.id, c.postID, c.userID, c.content, c.likes, c.createdAt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to seed comment %s: %v\n", c.id, err)
		}
	}
}
