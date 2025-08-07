Blog Routes


Public Routes

⦁	List Pots  - GET/blog

    Input: Not required
    	
    Response: {
        "limit": 10,
        "page": 1,
        "posts": [
            {
                "id": "689457b56e2cae04a9ace74d",
                "author_id": "6893544d594f56c731efd47d",
                "title": "Psychology",
                "content": "Psychology is the scientific study of the mind and behavior, exploring how people think, feel, and act. It helps us understand mental processes, emotions, and social interactions.",
                "tags": [
                    "mind",
                    "science",
                    "psycholgy"
                ],
                "view_count": 0,
                "likes": 0,
                "dislikes": 0,
                "comment_count": 0,
                "created_at": "2025-08-07T07:37:25.509Z",
                "updated_at": "2025-08-07T07:37:25.509Z"
            },
            {
                "id": "68940bef125e59b0b4f75335",
                "author_id": "6893544d594f56c731efd47d",
                "title": "Brain",
                "content": "The brain is the central organ of the human nervous system, responsible for processing information and controlling bodily functions. It enables thought, memory, emotion, and coordination, making it essential for all aspects of human life.",
                "tags": [
                    "mind",
                    "science",
                    "psycholgy"
                ],
                "view_count": 0,
                "likes": 1,
                "dislikes": 1,
                "comment_count": 0,
                "created_at": "2025-08-07T02:14:07.573Z",
                "updated_at": "2025-08-07T02:14:07.573Z"
            }
        ],
        "total": 2
    }

Optional Query Parameters:
 
  title with value string : Performs a case-insensitive "contains" search on the post title.
  
  author with value string: Finds posts by an author's full name (case-insensitive).
  
  tag with value string: Filters posts by one or more comma-separated tags.
  
  startDate with value date(YYYY-MM-DD): Returns posts created on or after this date.
   
  endDate with value date(YYYY-MM-DD): Returns posts created on or before this date.
  
  sortBy with value string: Specifies the sort order. Options: popularity, date_asc, date_desc. Default is date_desc.
  
  minPopularity with value integer: Returns posts with a likes count greater than or equal to this value.
  
  maxPopularity with value integer: Returns posts with a likes count less than or equal to this value.
  
  page with value integer: The page number for pagination. Default is 1. 
  
  limit with value integer: The number of posts to return per page. Default is 10.

⦁	Get Post By ID - GET/blog/:id

    Input: Not requires
    
    Response: {
        "id": "689457b56e2cae04a9ace74d",
        "author_id": "6893544d594f56c731efd47d",
        "title": "Psychology",
        "content": "Psychology is the scientific study of the mind and behavior, exploring how people think, feel, and act. It helps us understand mental processes, emotions, and social interactions.",
        "tags": [
            "mind",
            "science",
            "psycholgy"
        ],
        "view_count": 0,
        "likes": 0,
        "dislikes": 0,
        "comment_count": 0,
        "created_at": "2025-08-07T07:37:25.509Z",
        "updated_at": "2025-08-07T07:37:25.509Z"
    }

Protected Routes

⦁	Create Post - POST/blog

    Input: {
        "title": "Psychology",
        "content": "Psychology is the scientific study of the mind and behavior, exploring how people think, feel, and act. It helps us understand mental processes, emotions, and social interactions.",
        "tags": ["mind", "science", "psychology"]
    }
    
    Response: {
        "id": "689457b56e2cae04a9ace74d",
        "author_id": "6893544d594f56c731efd47d",
        "title": "Psychology",
        "content": "Psychology is the scientific study of the mind and behavior, exploring how people think, feel, and act. It helps us understand mental processes, emotions, and social interactions.",
        "tags": [
            "mind",
            "science",
            "psycholgy"
        ],
        "view_count": 0,
        "likes": 0,
        "dislikes": 0,
        "comment_count": 0,
        "created_at": "2025-08-07T00:37:25.5096018-07:00",
        "updated_at": "2025-08-07T00:37:25.5096018-07:00"
    }

⦁	Update Post - PUT/blog/:id

    Input: {
        "title": "Programming Language ",
        "content": "A programming language is a formal set of instructions used to communicate with computers and create software applications. It allows developers to write code that a machine can interpret and execute to perform specific tasks.",
        "tags": ["Go", "Python", "Java"]
    }
    
    
    Response: {
        "id": "68936643594f56c731efd482",
        "author_id": "68935bee594f56c731efd47f",
        "title": "Programming Language ",
        "content": "A programming language is a formal set of instructions used to communicate with computers and create software applications. It allows developers to write code that a machine can interpret and execute to perform specific tasks.",
        "tags": [
            "Go",
            "Python",
            "Java"
        ],
        "view_count": 0,
        "likes": 0,
        "dislikes": 0,
        "comment_count": 2,
        "created_at": "2025-08-06T14:27:15.66Z",
        "updated_at": "2025-08-07T00:41:30.4120417-07:00"
    }

⦁	Delete Post - DELETE/blog/:id

    Input: Not required
    
    Response: {
        "message": "Post deleted successfully"
    }

⦁	Like Post - POST/blog/:id/like

    Input: Not required
    
    Response: {
        "message": "Post liked successfully"
    }

⦁	Dislike Post - POST/blog/:id/dislike

    Input: Not required
    
    Response: {
        "message": "Post disliked successfully"
    }

⦁	Create Comment - POST/blog/:id/comments

    Input: {
        "content": "This is the first comment from the user"
    }
    
    Response: {
        "id": "68945ac66e2cae04a9ace74e",
        "blog_id": "689457b56e2cae04a9ace74d",
        "author_id": "68935bee594f56c731efd47f",
        "content": "This is the first comment from the user",
        "created_at": "2025-08-07T00:50:30.8968678-07:00"
    }
