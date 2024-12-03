package main

import (
	"SocialMediaApp/internal/store"
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func Seed(db *sql.DB, store store.Storage, flags ...string) {
	ctx := context.Background()

	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("Error creating user:", err)
			return
		}
	}

	_ = tx.Commit()

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment:", err)
			return
		}
	}

	for _, flag := range flags {
		if flag == "follow" {
			users, err := getAllUsers(db, context.Background())
			if err != nil {
				log.Println("Error getting all the users:", err)
				return
			}

			for _, user := range users {
				start := rand.Intn(len(users) - 101)

				for i := start; i < start+100; i++ {
					follower := users[i]
					if follower.ID != user.ID {
						err := store.Followers.Follow(context.Background(), follower.ID, user.ID)
						if err != nil {
							log.Println("Error following:", err)
							return
						}
					}
				}
			}
		}
	}

	log.Println("Seeding complete")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		timeNow := time.Now().Unix()
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", timeNow),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", timeNow) + "@example.com",
		}
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		posts[i] = &store.Post{
			Content: contents[rand.Intn(len(contents))],
			Title:   titles[rand.Intn(len(titles))],
			UserID:  user.ID,
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}

	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	comments := make([]*store.Comment, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		post := posts[rand.Intn(len(posts))]
		comments[i] = &store.Comment{
			PostID:  post.ID,
			UserID:  user.ID,
			Content: hardcodedComments[rand.Intn(len(hardcodedComments))],
			User:    *user,
		}
	}

	return comments
}

func getAllUsers(db *sql.DB, ctx context.Context) ([]*store.User, error) {
	query := `
		SELECT id
		FROM users 
	`

	rows, err := db.QueryContext(
		ctx,
		query,
	)

	if err != nil {
		log.Fatal("Error getting users from DB for seeding")
		return nil, err
	}

	var users []*store.User

	for rows.Next() {
		user := store.User{}
		err := rows.Scan(&user.ID)
		if err != nil {
			log.Fatal("Error reading a user from DB for seeding")
			return nil, err
		}

		users = append(users, &user)
	}

	return users, err
}

var usernames = []string{
	"James", "Mary", "John", "Patricia", "Robert", "Jennifer", "Michael", "Linda",
	"William", "Elizabeth", "David", "Barbara", "Richard", "Susan", "Joseph", "Jessica",
	"Thomas", "Sarah", "Charles", "Karen", "Christopher", "Nancy", "Daniel", "Lisa",
	"Matthew", "Betty", "Anthony", "Margaret", "Donald", "Sandra", "Mark", "Ashley",
	"Paul", "Dorothy", "Steven", "Kimberly", "Andrew", "Emily", "Kenneth", "Donna",
	"Joshua", "Michelle", "George", "Carol", "Kevin", "Amanda", "Brian", "Melissa",
	"Edward", "Deborah", "Ronald", "Stephanie", "Timothy", "Rebecca", "Jason", "Laura",
	"Jeffrey", "Sharon", "Ryan", "Cynthia", "Jacob", "Kathleen", "Gary", "Amy",
	"Nicholas", "Shirley", "Eric", "Angela", "Jonathan", "Helen", "Stephen", "Anna",
	"Larry", "Brenda", "Justin", "Pamela", "Scott", "Nicole", "Brandon", "Emma",
	"Frank", "Samantha", "Benjamin", "Katherine", "Gregory", "Christine", "Samuel", "Debra",
	"Raymond", "Rachel", "Patrick", "Catherine", "Alexander", "Carolyn", "Jack", "Janet",
	"Dennis", "Ruth", "Jerry", "Maria", "Tyler", "Heather",
}

var titles = []string{
	"Chasing sunsets", "Coffee first, always", "Weekend vibes", "Adventure awaits",
	"Smiles for days", "Nature's masterpiece", "City lights calling", "Dream big, hustle hard",
	"Exploring new horizons", "Golden hour moments", "Capturing the little things", "Making memories",
	"Everyday magic", "The calm before the storm", "Life on the go", "A moment in time",
	"Good food, good mood", "Sundays are for relaxing", "Work hard, play harder", "Simple joys",
	"Life is an adventure", "New beginnings", "The beauty of the ordinary", "Through my lens",
	"Finding balance", "Keep it simple", "Forever exploring", "Peaceful moments", "The journey matters",
	"Small steps every day", "Life's colorful moments", "Close to nature", "Every detail counts",
	"Lost in the moment", "Living the dream", "The power of perspective", "Fresh air and freedom",
	"Unforgettable journeys", "Happiness is homemade", "Time stands still", "Paving my path",
	"The art of being present", "Every day is a new story", "Moments worth sharing", "Beauty in simplicity",
	"Learning and growing", "Captured moments", "On top of the world", "In the heart of the city",
	"Grateful for today", "Savoring the view", "Wanderlust calling", "A life well lived",
	"Never stop dreaming", "Enjoying the little things", "Making it count", "Living my best life",
	"Step by step", "Collecting moments, not things", "Seeing the world differently", "Breathe and believe",
	"Every day is a gift", "Follow the light", "Finding my happy place", "Moments of gratitude",
	"Life in focus", "Pause and reflect", "Unfolding the story", "Making it happen",
	"One day at a time", "Inspired by nature", "Through the seasons", "The world is my canvas",
	"Joy in every step", "Reflections of today", "Taking a leap", "Chasing my dreams",
	"Creating something beautiful", "Fresh starts", "Hello, new day", "The charm of the unknown",
	"Simply living", "Moments to cherish", "Life is what you make it", "Going with the flow",
	"A step closer to the dream", "Today is the day", "In the moment", "Wandering through life",
	"A story worth sharing", "Embracing the now", "Focusing on the good", "Every corner tells a story",
	"Daydream believer", "Simplicity is key", "Turning dreams into plans", "A life of meaning",
}

var contents = []string{
	"Spent the day exploring the countryside and soaking in the fresh air.",
	"There's nothing better than a good book and a cozy corner.",
	"Captured this beautiful view during my morning hike.",
	"Cooked up a storm in the kitchen today – recipe coming soon!",
	"Reflecting on the simple joys of life that make everything worthwhile.",
	"Finally ticked this destination off my bucket list.",
	"Learning to appreciate the beauty of everyday moments.",
	"Here's to new beginnings and endless possibilities.",
	"Sometimes, the best days are the unplanned ones.",
	"Challenging myself to step out of my comfort zone today.",
	"Found this hidden gem while wandering through the city streets.",
	"A little progress every day adds up to big results.",
	"Throwback to a day filled with sunshine and laughter.",
	"Celebrating small wins that lead to big changes.",
	"Exploring the unknown and loving every second of it.",
	"Creating memories that I'll treasure forever.",
	"Today was all about reconnecting with nature.",
	"Letting go of the past and focusing on the future.",
	"Finally took some time to relax and recharge.",
	"Taking a moment to appreciate how far I've come.",
	"Started the day with a sunrise that took my breath away.",
	"Feeling inspired by the colors of the season.",
	"Spent time catching up with old friends and making new memories.",
	"Some days are meant for adventure, others for reflection.",
	"Found joy in the little things today – like a cup of coffee.",
	"Life is a journey, and I'm enjoying every step.",
	"Took a long walk and let my thoughts wander freely.",
	"Today was a reminder that hard work always pays off.",
	"Rediscovering the beauty of my own backyard.",
	"Feeling grateful for the people who make life meaningful.",
	"Every day is a chance to write a new story.",
	"Chased my dreams and found a little magic along the way.",
	"Captured this photo just as the sun was setting.",
	"Celebrating the beauty of imperfection.",
	"Learning to embrace change and all the opportunities it brings.",
	"Nothing beats a quiet evening spent by the fireplace.",
	"Spontaneous road trips always lead to the best stories.",
	"Tried something new today and loved the experience.",
	"Life feels better when you're surrounded by love and laughter.",
	"Discovering that the best views often come after the hardest climbs.",
	"Started a new project today – can't wait to share more!",
	"Feeling at peace with where I am and where I'm going.",
	"Some moments are too beautiful not to share.",
	"Made a small change today that made a big difference.",
	"Life's better when you stop to smell the flowers.",
	"Sometimes, all you need is a little fresh air and perspective.",
	"Taking time to pause and enjoy the journey.",
	"Today was all about self-care and self-love.",
	"Found beauty in the simplest things today.",
	"Living in the present and making it count.",
	"Every day is an opportunity to grow and learn.",
	"Started the day with gratitude and ended it with joy.",
	"Revisiting old places and discovering new things.",
	"Feeling recharged after a weekend spent outdoors.",
	"Sometimes, the best moments are the quiet ones.",
	"Spent the day doing what I love – and it feels amazing.",
	"Grateful for the moments that take my breath away.",
	"Living the life I've always dreamed of – one day at a time.",
	"Taking on new challenges and embracing the journey.",
	"Learning to see the world through fresh eyes.",
	"Today was a mix of work and play – the perfect balance.",
	"Spent the evening watching the stars and reflecting on life.",
	"Finding happiness in the most unexpected places.",
	"Every new day brings endless possibilities.",
	"Rediscovered an old passion today – and it felt great.",
	"Making the most of every moment – big or small.",
	"Creating memories that will last a lifetime.",
	"Feeling inspired by the beauty of the natural world.",
	"Nothing beats the feeling of accomplishing a goal.",
	"Spent the day exploring, learning, and growing.",
	"Taking time to slow down and enjoy the little things.",
	"Sometimes, it's the smallest changes that make the biggest impact.",
	"Found peace in the quiet moments today.",
	"Making space for new opportunities and adventures.",
	"Spent the day doing something I love with people I cherish.",
	"Celebrating the joy of creativity and self-expression.",
	"Found strength in unexpected places today.",
	"Grateful for the journey and the lessons along the way.",
	"Every moment is a chance to start fresh.",
	"Rediscovering the magic in everyday life.",
	"Spent the day chasing dreams and capturing memories.",
	"Nothing compares to the feeling of being truly present.",
	"Learning to let go and enjoy the ride.",
	"Finding beauty in every corner of the world.",
	"Today was a reminder of the power of persistence.",
	"Creating something meaningful, one step at a time.",
	"Spent the day reflecting on how far I've come.",
	"Taking time to nurture my passions and my soul.",
	"Exploring new possibilities and pushing boundaries.",
	"Today was a reminder to appreciate the journey.",
	"Found joy in the unexpected moments of the day.",
	"Living for the moments that make me smile.",
	"Spent the day surrounded by inspiration and creativity.",
	"Today was all about making dreams a reality.",
}

var tags = []string{
	"#NaturePhotography",
	"#TravelGoals",
	"#SelfCare",
	"#Foodie",
	"#FitnessJourney",
	"#Throwback",
	"#WeekendVibes",
	"#SunsetLovers",
	"#Wanderlust",
	"#DailyMotivation",
	"#ArtisticExpressions",
	"#CozyCorner",
	"#AdventureAwaits",
	"#CityLife",
	"#MountainEscape",
	"#DreamBig",
	"#GratefulHeart",
	"#NewBeginnings",
	"#BeachDay",
	"#MorningHike",
	"#BookLover",
	"#CreativeLife",
	"#RoadTrip",
	"#CalmMind",
	"#LoveAndLight",
	"#BeautifulViews",
	"#SmallSteps",
	"#SunnyDays",
	"#SimpleJoys",
	"#FamilyTime",
	"#NatureEscape",
	"#FreshPerspective",
	"#PeacefulMoments",
	"#InspiredLiving",
	"#HiddenGems",
	"#LoveForTravel",
	"#SelfDiscovery",
	"#MindfulnessMatters",
	"#UrbanAdventures",
	"#SlowLiving",
	"#CelebrateLife",
	"#BucketList",
	"#PositiveVibesOnly",
	"#ColorfulWorld",
	"#OutdoorFun",
	"#HikingAdventures",
	"#CherishMoments",
	"#NewHorizons",
	"#NatureLover",
	"#FocusOnTheGood",
	"#BigDreams",
	"#HealthyHabits",
	"#SeizeTheDay",
	"#GoodVibes",
	"#TBT",
	"#LifeLessons",
	"#JoyfulLiving",
	"#Tranquility",
	"#EpicAdventures",
	"#SunshineStateOfMind",
	"#AdventureLife",
	"#SoulfulLiving",
	"#FindYourPassion",
	"#UnforgettableJourneys",
	"#WorkHardPlayHard",
	"#SkylineViews",
	"#ChasingDreams",
	"#SereneScenes",
	"#PicturePerfect",
	"#LifeIsBeautiful",
	"#NatureInspired",
	"#FeelTheVibes",
	"#MomentsMatter",
	"#BoundlessEnergy",
	"#FreshAirAndFreedom",
	"#StayCurious",
	"#OutdoorAdventures",
	"#UrbanExploration",
	"#MemoriesMade",
	"#VibrantColors",
	"#HappinessIsHere",
	"#BeyondTheHorizon",
	"#DailyInspiration",
	"#RelaxAndRecharge",
	"#InfinitePossibilities",
	"#BeachVibes",
	"#PeaceOfMind",
	"#LifeUnfiltered",
	"#AdventuresInNature",
	"#CreateAndInspire",
	"#ExploreMore",
	"#DreamItDoIt",
	"#BrightFuture",
	"#NatureTherapy",
	"#CalmAndCollected",
	"#ChasingSunsets",
	"#LifeWellLived",
	"#PurposeAndPassion",
	"#StayGrounded",
	"#FindYourJoy",
	"#EndlessAdventures",
	"#CelebrateEveryDay",
	"#NatureAtItsBest",
}

var hardcodedComments = []string{
	"Absolutely love this!",
	"Stunning shot!",
	"Such a beautiful moment captured.",
	"Incredible vibes here!",
	"This is so inspiring!",
	"Great perspective in this photo.",
	"Can't wait to visit this place.",
	"Wow, just breathtaking!",
	"Thanks for sharing this beauty!",
	"This looks amazing, well done!",
	"Perfectly captured!",
	"I love everything about this post!",
	"Your content always inspires me.",
	"Amazing view!",
	"Such a creative shot.",
	"This just made my day!",
	"Pure perfection!",
	"Adding this to my bucket list!",
	"This is absolutely stunning.",
	"Your work is always incredible.",
	"Keep up the amazing content!",
	"Such a mood in this post.",
	"I can feel the vibe from here.",
	"Totally relatable!",
	"This is my kind of aesthetic.",
	"Fantastic colors in this one!",
	"Really resonates with me.",
	"This makes me so happy!",
	"Unbelievably beautiful.",
	"Such a peaceful scene.",
	"Thanks for brightening my day!",
	"Your posts are always on point.",
	"This is pure art!",
	"Really feeling this vibe.",
	"You've outdone yourself!",
	"Wow, this place looks magical.",
	"Such an iconic shot!",
	"This reminds me of my travels.",
	"Can't get enough of your posts!",
	"Always a pleasure to see your work.",
	"Absolutely flawless!",
	"This post just uplifted my mood.",
	"Such a meaningful capture.",
	"This inspires me to explore more.",
	"Amazing composition!",
	"Perfectly timed shot.",
	"Such a joyful post!",
	"This place looks surreal.",
	"I'm loving the colors in this!",
	"This is truly amazing.",
	"Such a great moment captured.",
	"Your photos always tell a story.",
	"This is picture-perfect!",
	"I'm in awe of this scene.",
	"This gives me so much peace.",
	"I can't stop staring at this!",
	"This caption is everything.",
	"Absolutely dreamy!",
	"This post is a masterpiece.",
	"I feel so inspired right now.",
	"This brings back so many memories.",
	"Such an iconic moment!",
	"I love the energy here!",
	"Your work keeps getting better.",
	"Always a fan of your style.",
	"This is pure gold.",
	"This speaks to my soul.",
	"Such a unique perspective!",
	"This place is now on my list!",
	"Such powerful imagery.",
	"I'm speechless!",
	"Your creativity knows no bounds.",
	"This made me smile!",
	"Such a peaceful vibe.",
	"This is the definition of beauty.",
	"Such a wonderful moment.",
	"Love the minimalistic vibe.",
	"How do you always capture this magic?",
	"This is so refreshing to see.",
	"Such positive energy in this post!",
	"Completely in love with this.",
	"Your work is so captivating.",
	"This is the best thing I've seen today.",
	"Such elegance in this post.",
	"I'm so glad I saw this!",
	"This is next-level amazing.",
	"I'm saving this for inspiration!",
	"This fills me with joy.",
	"This is truly breathtaking.",
	"Such a heartfelt capture.",
	"Absolutely incredible scene.",
	"This just made my evening!",
	"Can't wait to see more like this.",
	"Such a great moment in time.",
	"This brings me so much peace.",
	"This is why I love following you.",
}
