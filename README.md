# Messenger
A simples Go messenger application created to improve Golang skills. </br>
In this application you can create users, authenticated with AWS cognito, search for other active users and exchange messages with them.

## Endpoints
There are open endpoints and endpoints that need authorization. Authorization is granted to confirmed users after login.

#### Open Endpoints
<lu>
	<li><b>POST /api/v0/users</b> -> Creates user. Expects body with email, nickName and password.</li>
	<li><b>POST /api/v0/users/confirmation</b> -> Confirms user. Expects body with email and code (received by email).</li>
	<li><b>POST /api/v0/users/login</b> -> Logs in user. Expects body with email and password. Returns token.</li>
</lu>

### Authorized only endpoints
To access these endpoints bearer token authporization is required.
<lu>
	<li><b>GET /api/v0/user</b> -> Get token user's information. </li>
	<li><b>GET /api/v0/users</b> -> List users (limit 20). Optional: parameter name to filter email and username by subquery.</li>
	<li><b>PUT /api/v0/users/password</b> -> Updates token user password. Expects body with email and new password.</li>
	<li><b>DELETE /api/v0/users</b> -> Deletes token user.</li>
	<li><b>GET /api/v0/messages/{username}</b> -> Lists messages (limit 20) between token user and username user.</li>
	<li><b>/api/v0/chat/{username}</b> -> Establishes websocket connection to send and receive messages between token user and username user. If username user is also connected, messages can be exchanged live. </li>
</lu>

## Comments and future improvements
Authorized endpoints are a bit redundant, authorization wise and user wise. I was looking for a way to handle all authorized connections in one place but couldn't find, but that's an improvement I'd work on. Also I needed a local users table to list and filter them, but creates some seemenly code redundancies.
Endpoints are a bit out of pattern, for my linking. For instance, an endpoint that gives a user informations should be "GET /users/{id}", but the user already have the authorization token and, for now, doesn't have access to other users, so it made sense to use just "GET /user" with bearer token authorization. This was a choice, I guess, I could have gone the other way.
Also, the messages and chat endpoints uses username. My first thought was to use Id, but for me a "/messages/{id}" endpoint would suggest getting a specif message, not messages exchanged between 2 users.</br>
More improvements:
<lu>
  <li>Update username endpoint</li>
  <li>List users pagination</li>
</lu>