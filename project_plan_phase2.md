### Project Enhancement Plan: Agent & Customer Features

#### 1. Authentication and Authorization

The first step is to build a secure authentication system. This will be the foundation for securing the agent and customer-specific routes.

**Technology Choice:** JSON Web Tokens (JWT) are a good choice for this API-based system.

**Milestones:**

*   **1.1: Implement JWT Generation:**
    *   Create a new `/login` endpoint.
    *   Users will `POST` their `email` and `password`.
    *   The system will validate the credentials against the `users` table.
    *   On success, it will generate a JWT containing the `user_id` and `role`.
*   **1.2: Create Authentication Middleware:**
    *   This middleware will protect the `/agent` and `/customer` routes.
    *   It will inspect the `Authorization` header for a valid JWT.
    *   It will decode the token to identify the user and their role.
    *   If the token is invalid or the role does not match the route's requirement, it will return a `401 Unauthorized` or `403 Forbidden` error.

---

#### 2. Agent Features

Agents are the core of the ticketing system. They need tools to manage the tickets assigned to them efficiently.

**Route Group:** `/agent` (All endpoints here will require `Agent` role authentication)

**Milestones:**

*   **2.1: View and Manage Assigned Tickets:**
    *   `GET /agent/tickets`: List all tickets where the `assignee_id` matches the authenticated agent's ID.
    *   `GET /agent/tickets/{id}`: View the details of a specific ticket, but only if assigned to them.
    *   `PUT /agent/tickets/{id}`: Update the `status` (e.g., 'In Progress', 'Resolved') or `priority` of an assigned ticket.
*   **2.2: Comment on Tickets:**
    *   `POST /agent/tickets/{id}/comments`: Add a comment to a ticket. The request body should allow specifying if the comment `is_internal` (visible only to other agents/admins) or public (visible to the customer).
*   **2.3: View Customer Information:**
    *   When viewing a ticket, the API response should include the requester's (customer's) details (`name`, `email`).

---

#### 3. Customer Features

Customers need a simple way to create, track, and communicate about their support requests.

**Route Group:** `/customer` (All endpoints here will require `Customer` role authentication)

**Milestones:**

*   **3.1: Manage Their Own Tickets:**
    *   `POST /customer/tickets`: Create a new ticket. The `requester_id` will be automatically set to the authenticated customer's ID.
    *   `GET /customer/tickets`: List all tickets where the `requester_id` matches the authenticated customer's ID.
    *   `GET /customer/tickets/{id}`: View the details of a specific ticket they created.
*   **3.2: Communication:**
    *   `POST /customer/tickets/{id}/comments`: Add a public comment to one of their own tickets. The API will ensure `is_internal` is always `false`.
*   **3.3: Ticket Closure:**
    *   `PUT /customer/tickets/{id}`: Allow a customer to change the status of their ticket to `Closed`.
