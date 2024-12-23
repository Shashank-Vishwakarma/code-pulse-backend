package models

import (
	"math"
	"time"
)

type Blog struct {
	ID              string    `json:"id" bson:"_id"`
	Title           string    `json:"title" bson:"title"`
	Body            string    `json:"body" bson:"body"`
	IsBlogPublished bool      `json:"isBlogPublished" bson:"isBlogPublished"`
	Slug            string    `json:"slug" bson:"slug"`
	Comments        []Comment `json:"comments" bson:"comments"`
	UpVotes         []string  `json:"upVotes" bson:"upVotes"`
	DownVotes       []string  `json:"downVotes" bson:"downVotes"`
	Score           float64   `json:"score" bson:"score"` // to track top 10 blogs
	AuthorID        string    `json:"authorId" bson:"authorId"`
	CreatedAt       time.Time `json:"createdAt" bson:"createdAt"`
}

func (b *Blog) CalculateBlogScore() float64 {
	// Calculate net votes
	netVotes := len(b.UpVotes) - len(b.DownVotes)

	// percentage of downvotes
	downvotesPercentage := float64(len(b.DownVotes)) / float64(len(b.UpVotes)+len(b.DownVotes))
	if downvotesPercentage > 0.2 {
		netVotes = 0
	}

	// Calculate vote score with logarithmic scaling to prevent vote count inflation
	var voteScore float64
	signedBit := math.Signbit(float64(netVotes))
	if signedBit {
		voteScore = math.Log1p(math.Abs(float64(netVotes)))
	} else {
		voteScore = 0
	}

	// Calculate recency score (blogs within last 30 days get higher scores)
	daysSinceCreation := time.Since(b.CreatedAt).Hours() / 24
	recencyScore := math.Max(0, 30-daysSinceCreation) / 30

	// Calculate engagement score
	commentScore := math.Log1p(float64(len(b.Comments)))

	// Combine scores with weights
	// Adjust these weights as needed
	totalScore := voteScore*0.5 + recencyScore*0.3 + commentScore*0.2

	// Normalize and round the score
	b.Score = math.Round(totalScore*100) / 100
	return b.Score
}
