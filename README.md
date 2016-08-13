A website and user system starter. Implemented with gin and Backbone.
Gowall is port of [Drywall](https://github.com/jedireza/drywall)

|            | Go                                            | Node.js                                         |
| ---------- | --------------------------------------------- | ----------------------------------------------- |
| Repository | here                                          | [Drywall](https://github.com/jedireza/drywall/) |
| Site       | [Gowall](http://im7mortal.github.io/gowall/)  | [Drywall](http://jedireza.github.io/drywall/)   |
| Demo       | [Gowall demo](https://go-wall.herokuapp.com/) | [Drywall demo](https://drywall.herokuapp.com/)  |

I cloned Drywall from [commit](https://github.com/jedireza/drywall/tree/ca35c06bf5100d2835da929bd1bd3c39ae441138)

## Technology

Server side, Gowall is built with the [gin](https://github.com/gin-gonic/gin)
framework. We're using [MongoDB](http://www.mongodb.org/) as a data store.

The front-end is built with [Backbone](http://backbonejs.org/).
You can use [Grunt](http://gruntjs.com/) for the asset pipeline.
Grunt's settings are located in Drywall's repository.

| On The Server      | On The Client  | Development |
| ------------------ | -------------- | ----------- |
| Gin                | Bootstrap      | Grunt       |
| html.template      | Backbone.js    |             |
| mgo                | jQuery         |             |
| markbates/goth     | Underscore.js  |             |
| gopkg.in/gomail.v2 | Font-Awesome   |             |
|                    | Moment.js      |             |


## Live demo

| Platform                       | Username | Password |
| ------------------------------ | -------- | -------- |
| https://go-wall.herokuapp.com/ | root     | h3r00t   |

__Note:__ The live demo has been modified so you cannot change the root user,
the root user's linked admin role or the root admin group. This was done in
order to keep the app ready to use at all times.

## Diference

I keer session in cookies.

## Requirements

You need [MongoDB](http://www.mongodb.org/downloads) installed and running.

We use [`mgo`](https://labix.org/mgo) as mongodb driver. If
you have issues with sasl [refer to this issue](https://github.com/go-mgo/mgo/issues/220#issuecomment-212605949).

We use `golang.org/x/crypto/bcrypt` for hashing
secrets.

## Installation

Exetuble file has to be located in the same directory where ~public~ is located


## Setup

First you need to setup your config file. It's located just inside. cmd/gowall/config.example.go
I decided that you know better how you want keep your config.
It support system environments.

Next, you need a few records in the database to start using the user system.

Run these commands on mongo via the terminal. __Obviously you should use your
email address.__

```js
use Gowall; // or your mongo db name if different
```

```js
db.admingroups.insert({ _id: 'root', name: 'Root' });
db.admins.insert({ name: {first: 'Root', last: 'Admin', full: 'Root Admin'}, groups: ['root'] });
var rootAdmin = db.admins.findOne();
db.users.save({ username: 'root', isActive: 'yes', email: 'your@email.addy', roles: {admin: rootAdmin._id} });
var rootUser = db.users.findOne();
rootAdmin.user = { id: rootUser._id, name: rootUser.username };
db.admins.save(rootAdmin);
```


## Running the app

```bash
$ ./gowall
```

Now just use the reset password feature to set a password.

 - Go to `http://localhost:3000/login/forgot/`
 - Submit your email address and wait a second.
 - Go check your email and get the reset link.
 - `http://localhost:3000/login/reset/:email/:token/`
 - Set a new password.

Login. Customize. Enjoy.

## Error handling

I use exception model for errors from std or third part libraries.
If mgo.Query.Find return err != mgo.ErrNotFound  I do panic
If mgo.ErrNotFound It's common behaviour except some cases when one object is nonsense if other object doesn't exist.
(I mean that case when necessary manual update of db)
In that case I am doing panic but with my own error description. That sysadmin could see it errors in log.
I like go error flow idea but if error has to be written in log i do panic.
I created func EXCEPTION(i interface{})  where you can specify your handler/ You can even remove panic from here.


## Philosophy

 - Create a website and user system.
 - Write code in a simple and consistent way.
 - Only create minor utilities or plugins to avoid repetitiveness.
 - Find and use good tools.
 - Use tools in their native/default behavior.


## Features

 - Basic front end web pages.
 - Contact page has form to email.
 - Login system with forgot password and reset password.
 - Signup and Login with Facebook, Twitter, GitHub, Google and Tumblr.
 - Optional email verification during signup flow.
 - User system with separate account and admin roles.
 - Admin groups with shared permission settings.
 - Administrator level permissions that override group permissions.
 - Global admin quick search component.


## Questions and contributing

Any issues or questions (no matter how basic), open an issue. Please take the
initiative to include basic debugging information like operating system
and relevant version details such as:

If you're changing something non-trivial, you may want to submit an issue
first.

## Thanks

Big thanks to [@jedireza](https://twitter.com/jedireza) for Drywall!!

## License

MIT

1. Can't do dynamic providers
2. Can't init check hostName
