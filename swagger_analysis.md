# Swagger API Analysis

## API Overview
- **Title**: WebApi$1.0
- **Version**: 1.0
- **Host**: mysource-gcdxbbb8akfge3du.northeurope-01.azurewebsites.net
- **Protocol**: HTTPS
- **Swagger Version**: 2.0

## Authentication
- **Type**: Bearer Token (JWT)
- **Header**: Authorization
- **Format**: "Bearer {token}"

## API Endpoints Summary

### Authentication Endpoints (`/auth/v1/`)
1. **POST /auth/v1/token** - Authenticate user with internal email and password
2. **GET /auth/v1/external-login-{provider}** - Redirect to External Provider for authentication
3. **GET /auth/v1/external-callback-{provider}** - Authenticate user with External account
4. **POST /auth/v1/google-token** - Authenticate user with Google token
5. **POST /auth/v1/signup** - Signup with email/password; sends verification code to email
6. **POST /auth/v1/verify-email** - Verify email with code and auto-login

### REST API Endpoints (`/rest/v1/`)

#### Category Management
- **POST /rest/v1/Category** - Create a new category

#### Conversation Management
- **POST /rest/v1/Conversation** - Create a new conversation

#### Escrow Management
- **POST /rest/v1/Escrow** - Create a new escrow
- **PUT /rest/v1/Escrow/{id}/status** - Update escrow status

#### File Management
- **POST /rest/v1/EscrowFile** - Create a new escrow file
- **GET /rest/v1/EscrowFile/{escrowFileId}** - Retrieves a single escrow file
- **POST /rest/v1/File** - Create a new file
- **GET /rest/v1/File/{fileId}** - Retrieves a single file
- **POST /rest/v1/FileGroup** - Create a new file group

#### Message Management
- **POST /rest/v1/Message** - Create a new message
- **PUT /rest/v1/Message/{id}/read** - Mark a single message as read
- **PUT /rest/v1/Message/bulk/read** - Mark multiple messages as read

#### Notification Management
- **POST /rest/v1/Notification** - Create a new notification
- **POST /rest/v1/Notification/register** - Register notification token
- **POST /rest/v1/Notification/test** - Test notification
- **POST /rest/v1/Notification/topic_test** - Test topic notification

#### Profile Management
- **POST /rest/v1/Profile** - Create/update profile
- **POST /rest/v1/Profile/avatar** - Upload profile avatar
- **GET /rest/v1/Profile/{storagePath}** - Get profile by storage path

#### Request Management
- **POST /rest/v1/Request** - Create a new request
- **POST /rest/v1/RequestResponse** - Create a new request response

#### Research Database
- **POST /rest/v1/ResearchDatabase** - Create a new research database entry

#### Resource Declaration
- **POST /rest/v1/ResourceDeclaration** - Create a new resource declaration
- **POST /rest/v1/ResourceDeclaration/bulk** - Bulk operations on resource declarations
- **GET /rest/v1/ResourceDeclaration/{id}** - Get specific resource declaration

#### User Management
- **POST /rest/v1/User** - Create a new user
- **GET /rest/v1/User/{id}** - Get specific user
- **POST /rest/v1/UserLocation** - Create user location

#### System
- **GET /rest/v1/version** - Get API version

### OData Endpoints (`/odata/`)
The API also provides OData endpoints for querying data with standard OData operations:
- `/odata/Category` - Category data with OData querying
- `/odata/Conversation` - Conversation data with OData querying
- `/odata/Escrow` - Escrow data with OData querying
- `/odata/EscrowFile` - Escrow file data with OData querying
- `/odata/File` - File data with OData querying
- `/odata/FileGroup` - File group data with OData querying
- `/odata/Message` - Message data with OData querying
- `/odata/Notification` - Notification data with OData querying
- `/odata/Profile` - Profile data with OData querying
- `/odata/Request` - Request data with OData querying
- `/odata/RequestResponse` - Request response data with OData querying
- `/odata/ResearchDatabase` - Research database data with OData querying
- `/odata/ResourceDeclaration` - Resource declaration data with OData querying
- `/odata/UserLocation` - User location data with OData querying
- `/odata/hotjobs` - Hot jobs data
- `/odata/$metadata` - OData metadata
- `/odata` - OData service document

## Key Data Models

### Authentication Models
- **LoginCommand**: Email and password for authentication
- **GoogleTokenLoginCommand**: Google token for authentication
- **SignupCommand**: User registration data
- **VerifyEmailCommand**: Email verification with code

### Business Models
- **CreateCategoryCommand**: Category creation with name, description, and icon
- **CreateConversationCommand**: Conversation creation with request/response IDs
- **CreateEscrowCommand**: Escrow creation with conversation, amount, description, and status
- **CreateFileGroupCommand**: File group creation with conversation and uploader IDs
- **CreateMessageCommand**: Message creation with content, conversation, type, and sender
- **CreateNotificationCommand**: Notification creation with user ID, title, and message
- **CreateRequestCommand**: Request creation
- **CreateUserDto**: User creation data
- **UpdateUserDto**: User update data with full name, active status, and role

## API Features
1. **Comprehensive Authentication**: Multiple authentication methods including internal, external providers, and Google
2. **File Management**: Complete file handling with escrow files, file groups, and profile avatars
3. **Messaging System**: Real-time messaging with read status tracking
4. **Notification System**: Push notifications with token registration and testing
5. **OData Support**: Full OData querying capabilities for all major entities
6. **User Management**: Complete user lifecycle management
7. **Request/Response System**: Structured request and response handling
8. **Research Database**: Research data management
9. **Resource Declaration**: Resource management and bulk operations

## Content Types Supported
The API supports extensive OData content types including:
- Various JSON formats with different metadata levels
- OData streaming options
- IEEE754 compatibility options
- XML and plain text formats

## Security
- JWT Bearer token authentication required for most endpoints
- API versioning support via query parameter and header
- Comprehensive input validation on all endpoints