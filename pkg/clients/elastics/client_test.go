package elastics

import (
	"bytes"
	"context"
	"testing"
)

func TestNewClient(t *testing.T) {

	c := NewClient(Config{
		Addresses: []string{"http://localhost:9200"},
	})

	mapping := `
 {
  "mappings": {
    "properties": {
      "title": {
        "type": "text",
        "analyzer": "ik_max_word",
        "fields": {
          "keyword": {
            "type": "keyword"
          }
        }
      },
      "status": {
        "type": "keyword"
      },
      "coverUrl": {
        "type": "keyword"
      },
      "url": {
        "type": "keyword"
      },
      "commodity": {
        "properties": {
          "brand": {
            "type": "text",
            "analyzer": "ik_max_word",
            "fields": {
              "keyword": {
                "type": "keyword"
              }
            }
          },
          "name": {
            "type": "text",
            "analyzer": "ik_max_word",
            "fields": {
              "keyword": {
                "type": "keyword"
              }
            }
          }
        }
      },
      "createdAt": {
        "type": "date",
        "format": "epoch_second"
      },
      "analysis": {
        "properties": {
          "segments": {
            "type": "nested",
            "properties": {
              "category": {
                "type": "text",
                "analyzer": "ik_max_word",
                "fields": {
                  "keyword": {
                    "type": "keyword"
                  }
                }
              },
              "description": {
                "type": "text",
                "analyzer": "ik_max_word",
                "fields": {
                  "keyword": {
                    "type": "keyword"
                  }
                }
              },
              "elements": {
                "properties": {
                  "personCount": {
                    "type": "text",
                    "fields": {
                      "keyword": {
                        "type": "keyword"
                      }
                    }
                  },
                  "personType": {
                    "type": "text",
                    "analyzer": "ik_max_word",
                    "fields": {
                      "keyword": {
                        "type": "keyword"
                      }
                    }
                  },
                  "voiceOverType": {
                    "type": "text",
                    "analyzer": "ik_max_word",
                    "fields": {
                      "keyword": {
                        "type": "keyword"
                      }
                    }
                  }
                }
              },
              "stage": {
                "properties": {
                  "name": {
                    "type": "text",
                    "analyzer": "ik_max_word",
                    "fields": {
                      "keyword": {
                        "type": "keyword"
                      }
                    }
                  },
                  "tags": {
                    "type": "text",
                    "analyzer": "ik_max_word",
                    "fields": {
                      "keyword": {
                        "type": "keyword"
                      }
                    }
                  }
                }
              },
              "timeStart": {
                "type": "float"
              },
              "timeEnd": {
                "type": "float"
              },
              "voiceOverScript": {
                "type": "text",
                "analyzer": "ik_max_word",
                "fields": {
                  "keyword": {
                    "type": "keyword"
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}

`

	do, err := c.GetElasticClient().Indices.Create("test").Raw(bytes.NewReader([]byte(mapping))).Do(context.Background())
	print(do, err)
}
