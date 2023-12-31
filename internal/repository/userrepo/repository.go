package userrepo

import (
	"context"
	"errors"
	"time"
	"vietcard-backend/internal/delivery/http/request"
	"vietcard-backend/internal/domain/entity"
	"vietcard-backend/internal/domain/interface/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	db      *mongo.Database
	colName string
}

func NewUserRepository(db *mongo.Database) repository.UserRepository {
	return &userRepository{
		db:      db,
		colName: "users",
	}
}

func (ur *userRepository) Create(user *entity.User) (string, error) {
	user.SetDefault()
	result, err := ur.db.Collection(ur.colName).InsertOne(context.TODO(), user)
	if err != nil {
		return "", err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("failed to get inserted ID")
	}

	return insertedID.Hex(), nil
}

func (ur *userRepository) GetByEmail(email *string) (*entity.User, error) {
	var user entity.User
	err := ur.db.Collection(ur.colName).FindOne(context.TODO(), bson.D{{Key: "email", Value: *email}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) GetByID(id *string) (*entity.User, error) {
	oID, err := primitive.ObjectIDFromHex(*id)
	if err != nil {
		return nil, err
	}
	var user entity.User
	err = ur.db.Collection(ur.colName).FindOne(context.TODO(), bson.D{{Key: "_id", Value: oID}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) UpdateUser(userID *string, req *request.UpdateUserRequest) (*entity.User, error) {
	uID, err := primitive.ObjectIDFromHex(*userID)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: uID}}
	update := bson.D{{Key: "$set", Value: *req}}
	option := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedUser entity.User
	err = ur.db.Collection(ur.colName).FindOneAndUpdate(context.TODO(), filter, update, option).Decode(&updatedUser)
	if err != nil {
		return nil, err
	}
	return &updatedUser, nil
}

func (ur *userRepository) UpdateUserXP(user *entity.User) error {
	filter := bson.D{{Key: "_id", Value: user.ID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "xp", Value: user.XP},
			{Key: "xp_to_level_up", Value: user.XPToLevelUp},
			{Key: "level", Value: user.Level},
			{Key: "streak", Value: user.Streak},
			{Key: "last_streak", Value: user.LastStreak},
		}},
	}
	_, err := ur.db.Collection(ur.colName).UpdateOne(context.TODO(), &filter, &update)
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) CreateFact(fact *entity.Fact) error {
	_, err := ur.db.Collection("fun_facts").InsertOne(context.TODO(), *fact)
	if err != nil {
		return err
	}

	return nil
}

func (ur *userRepository) GetFact() (*entity.Fact, error) {
	pipeline := bson.A{
		bson.D{{"$sample", bson.D{{"size", 1}}}},
	}

	opts := options.Aggregate().SetMaxTime(2 * time.Second) // You can adjust the timeout as needed

	cursor, err := ur.db.Collection("fun_facts").Aggregate(context.TODO(), pipeline, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	if !cursor.Next(context.TODO()) {
		return nil, errors.New("no facts found")
	}

	var result entity.Fact
	err = cursor.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
