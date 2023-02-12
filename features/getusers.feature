Feature: Users API (Get user)

    Scenario: GET a specific user
        Given I have these users in the database:
            """
            [
                {
                    "id": "ca99d09c-953a-4fe5-9b0a-51b3d40c01f7",
                    "name": "John Doe",
                    "email": "johndoe@postoffice.com",
                    "password": "bcryptesPassword"
                },
                {
                    "id": "aa99d09c-953a-4fe5-9b0a-51b3d40c01f7",
                    "name": "John Doe1",
                    "email": "johndoe1@postoffice.com",
                    "password": "bcryptesPassword123"
                }
            ]
            """
        When I GET "/users/ca99d09c-953a-4fe5-9b0a-51b3d40c01f7"
        Then I should receive the following model response with status "200":
            """
                {
                    "id": "ca99d09c-953a-4fe5-9b0a-51b3d40c01f7",
                    "name": "John Doe",
                    "email": "johndoe@postoffice.com"
                }
            """

    Scenario: GET a non existing user
        Given I have these users in the database:
            """
            [
                {
                    "id": "ca99d09c-953a-4fe5-9b0a-51b3d40c01f7",
                    "name": "John Doe",
                    "email": "johndoe@postoffice.com",
                    "password": "bcryptesPassword"
                },
                {
                    "id": "aa99d09c-953a-4fe5-9b0a-51b3d40c01f7",
                    "name": "John Doe1",
                    "email": "johndoe1@postoffice.com",
                    "password": "bcryptesPassword123"
                }
            ]
            """
        When I GET "/users/xxxxxxxx-953a-4fe5-9b0a-51b3d40c01f7"
        Then the HTTP status code should be "404"