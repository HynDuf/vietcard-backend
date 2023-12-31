package deckrepo

import (
	"context"
	"vietcard-backend/internal/delivery/http/request"
	"vietcard-backend/internal/domain/entity"
	"vietcard-backend/internal/domain/interface/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type deckRepository struct {
	db      *mongo.Database
	colName string
}

func NewDeckRepository(db *mongo.Database) repository.DeckRepository {
	return &deckRepository{
		db:      db,
		colName: "decks",
	}
}

func (dr *deckRepository) CreateDeck(deck *entity.Deck) (*entity.Deck, error) {
	deck.SetDefault()
	result, err := dr.db.Collection(dr.colName).InsertOne(context.TODO(), deck)
	if err != nil {
		return nil, err
	}
	deck.ID = result.InsertedID.(primitive.ObjectID)
	return deck, nil
}

func (dr *deckRepository) GetDeckByID(id *string) (*entity.Deck, error) {
	oID, err := primitive.ObjectIDFromHex(*id)
	if err != nil {
		return nil, err
	}
	var deck entity.Deck
	err = dr.db.Collection(dr.colName).FindOne(context.TODO(), bson.D{{Key: "_id", Value: oID}}).Decode(&deck)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &deck, nil
}

func (dr *deckRepository) UpdateDeck(deckID *string, req *request.UpdateDeckRequest) (*entity.Deck, error) {
	dID, err := primitive.ObjectIDFromHex(*deckID)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: dID}}
	update := bson.D{{Key: "$set", Value: *req}}
	option := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedDeck entity.Deck
	err = dr.db.Collection(dr.colName).FindOneAndUpdate(context.TODO(), filter, update, option).Decode(&updatedDeck)
	if err != nil {
		return nil, err
	}
	return &updatedDeck, nil
}

func (dr *deckRepository) GetCardsAllDecksOfUser(userID *string) (*[]entity.DeckWithCards, error) {
	ctx := context.TODO()

	uID, err := primitive.ObjectIDFromHex(*userID)
	if err != nil {
		return nil, err
	}

	// Open an aggregation cursor
	coll := dr.db.Collection(dr.colName)
	cursor, err := coll.Aggregate(ctx, bson.A{
		bson.D{{"$match", bson.D{{"user_id", uID}}}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "cards"},
					{"localField", "_id"},
					{"foreignField", "deck_id"},
					{"as", "cards"},
				},
			},
		},
		bson.D{
			{"$set",
				bson.D{
					{"cards",
						bson.D{
							{"$sortArray",
								bson.D{
									{"input", "$cards"},
									{"sortBy", bson.D{{"created_at", 1}}},
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var deckWithCards []entity.DeckWithCards
	if err = cursor.All(context.TODO(), &deckWithCards); err != nil {
		return nil, err
	}
	return &deckWithCards, nil
}

func (dr *deckRepository) GetCardsAllDecks(userID *string) (*[]entity.DeckWithCards, error) {
	// Requires the MongoDB Go Driver
	// https://go.mongodb.org/mongo-driver
	ctx := context.TODO()

	uID, err := primitive.ObjectIDFromHex(*userID)
	if err != nil {
		return nil, err
	}

	// Open an aggregation cursor
	coll := dr.db.Collection(dr.colName)
	cursor, err := coll.Aggregate(ctx, bson.A{
		bson.D{
			{"$match",
				bson.D{
					{"$or",
						bson.A{
							bson.D{{"user_id", uID}},
							bson.D{{"is_public", true}},
						},
					},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "cards"},
					{"localField", "_id"},
					{"foreignField", "deck_id"},
					{"as", "cards"},
				},
			},
		},
		bson.D{
			{"$set",
				bson.D{
					{"cards",
						bson.D{
							{"$sortArray",
								bson.D{
									{"input", "$cards"},
									{"sortBy", bson.D{{"created_at", 1}}},
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var deckWithCards []entity.DeckWithCards
	if err = cursor.All(context.TODO(), &deckWithCards); err != nil {
		return nil, err
	}
	return &deckWithCards, nil
}

func (dr *deckRepository) DeleteDeck(deckID *string) error {
	dID, err := primitive.ObjectIDFromHex(*deckID)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: dID}}
	_, err = dr.db.Collection(dr.colName).DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	filter = bson.D{{Key: "deck_id", Value: dID}}
	_, err = dr.db.Collection("cards").DeleteMany(context.TODO(), filter)
	if err != nil {
		return err
	}

	return nil
}

func (dr *deckRepository) GetDeckWithCards(deckID *string) (*entity.DeckWithCards, error) {
	ctx := context.TODO()

	dID, err := primitive.ObjectIDFromHex(*deckID)
	if err != nil {
		return nil, err
	}

	// Open an aggregation cursor
	coll := dr.db.Collection(dr.colName)
	cursor, err := coll.Aggregate(ctx, bson.A{
		bson.D{{"$match", bson.D{{"_id", dID}}}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "cards"},
					{"localField", "_id"},
					{"foreignField", "deck_id"},
					{"as", "cards"},
				},
			},
		},
		bson.D{
			{"$set",
				bson.D{
					{"cards",
						bson.D{
							{"$sortArray",
								bson.D{
									{"input", "$cards"},
									{"sortBy", bson.D{{"created_at", 1}}},
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var deckWithCardsList []entity.DeckWithCards
	if err = cursor.All(context.TODO(), &deckWithCardsList); err != nil {
		return nil, err
	}

	// Take the first element if available
	if len(deckWithCardsList) > 0 {
		return &deckWithCardsList[0], nil
	}

	// Handle case where no documents are found
	return nil, mongo.ErrNoDocuments
}
