package controllers

import (
	"context"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/BRAVO68WEB/go-rest-api/config"
	"github.com/BRAVO68WEB/go-rest-api/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllGists(c *fiber.Ctx) error {
    gistCollection := config.MI.DB.Collection("gists")
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

    var catchphrases []models.Catchphrase

    filter := bson.M{}
    findOptions := options.Find()

    if s := c.Query("s"); s != "" {
        filter = bson.M{
            "$or": []bson.M{
                {
                    "gistTitle": bson.M{
                        "$regex": primitive.Regex{
                            Pattern: s,
                            Options: "i",
                        },
                    },
                },
                {
                    "gistTopic": bson.M{
                        "$regex": primitive.Regex{
                            Pattern: s,
                            Options: "i",
                        },
                    },
                },
            },
        }
    }

    page, _ := strconv.Atoi(c.Query("page", "1"))
    limitVal, _ := strconv.Atoi(c.Query("limit", "10"))
    var limit int64 = int64(limitVal)

    total, _ := gistCollection.CountDocuments(ctx, filter)

    findOptions.SetSkip((int64(page) - 1) * limit)
    findOptions.SetLimit(limit)

    cursor, err := gistCollection.Find(ctx, filter, findOptions)
    defer cursor.Close(ctx)

    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "success": false,
            "message": "Gist Not found",
            "error":   err,
        })
    }

    for cursor.Next(ctx) {
        var catchphrase models.Catchphrase
        cursor.Decode(&catchphrase)
        catchphrases = append(catchphrases, catchphrase)
    }

    last := math.Ceil(float64(total / limit))
    if last < 1 && total > 0 {
        last = 1
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "data":      catchphrases,
        "total":     total,
        "page":      page,
        "last_page": last,
        "limit":     limit,
    })
}

func GetGist(c *fiber.Ctx) error {
    gistCollection := config.MI.DB.Collection("gists")
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

    var catchphrase models.Catchphrase
    objId, err := primitive.ObjectIDFromHex(c.Params("id"))
    findResult := gistCollection.FindOne(ctx, bson.M{"_id": objId})
    if err := findResult.Err(); err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "success": false,
            "message": "Gist Not found",
            "error":   err,
        })
    }

    err = findResult.Decode(&catchphrase)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "success": false,
            "message": "Gist Not found",
            "error":   err,
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "data":    catchphrase,
        "success": true,
    })
}

func AddGist(c *fiber.Ctx) error {
    gistCollection := config.MI.DB.Collection("gists")
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    catchphrase := new(models.Catchphrase)

    if err := c.BodyParser(catchphrase); err != nil {
        log.Println(err)
        return c.Status(400).JSON(fiber.Map{
            "success": false,
            "message": "Failed to parse body",
            "error":   err,
        })
    }

    result, err := gistCollection.InsertOne(ctx, catchphrase)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "success": false,
            "message": "Gist failed to insert",
            "error":   err,
        })
    }
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "data":    result,
        "success": true,
        "message": "Gist inserted successfully",
    })

}

func UpdateGist(c *fiber.Ctx) error {
    gistCollection := config.MI.DB.Collection("gists")
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    catchphrase := new(models.Catchphrase)

    if err := c.BodyParser(catchphrase); err != nil {
        log.Println(err)
        return c.Status(400).JSON(fiber.Map{
            "success": false,
            "message": "Failed to parse body",
            "error":   err,
        })
    }

    objId, err := primitive.ObjectIDFromHex(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "success": false,
            "message": "Gists not found",
            "error":   err,
        })
    }

    update := bson.M{
        "$set": catchphrase,
    }
    _, err = gistCollection.UpdateOne(ctx, bson.M{"_id": objId}, update)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "success": false,
            "message": "Gists failed to update",
            "error":   err.Error(),
        })
    }
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "success": true,
        "message": "Gists updated successfully",
    })
}

func DeleteGist(c *fiber.Ctx) error {
    gistCollection := config.MI.DB.Collection("gists")
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

    objId, err := primitive.ObjectIDFromHex(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "success": false,
            "message": "Gists not found",
            "error":   err,
        })
    }
    _, err = gistCollection.DeleteOne(ctx, bson.M{"_id": objId})
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "success": false,
            "message": "Gists failed to delete",
            "error":   err,
        })
    }
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "success": true,
        "message": "Gists deleted successfully",
    })
}