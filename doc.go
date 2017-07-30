/*
Package libplumraw is a client library for direct interactions with PlumLife and the Lightpad smart switches.

Summary

libplumraw does nothing more than turn the Plum REST and other APIs from HTTP
calls to function calls. It gives you an easy way to represent each type of
object you can manage via the Plum Web API and via the Lightpads themselves, as
well as listen for changes from a specific lightpad or listen for the lightpad
announcement heartbeats.

For higher level interactions, use `libplum` instead of `libplumraw`.

Use

To issue calls out to the website (to get house, room, etc. configs), first get
a `WebConnection` by calling `NewWebConnection()` with a `WebConnectionConfig`
(`Email` and `Password` are required, `PlumAPIHost` is optional). Use the
returned connection object to call out to the website.

There are three places from which you get information about Lightpads.

    * the general config comes from the Plum web service
    ** use `WebConnection.GetLightpad()` to fetch this data
    * the IP and Port come from the Heartbeat broadcast
    ** use `DefaultLightpadHeartbeat{}.Listen() to receive these messages
    * live changes to state come from a stream the lightpad itself produces
    ** use `Lightpad.Stream()` to get these updates

To interact with a Lightpad, you need to get its general config (including, for
example, the `Room` in which it exists) from the web. You then must merge that
with the IP and Port that come from the heartbeat broadcast. Finally, you should
listen to the update stream to adjust your internal data structures to represent
the current state of the switch. It is recomemended that you store this state
locally to avoid waiting 5 minutes for the heartbeat before being able to
interact with a Lightpad switch.

Testing code that uses libplumraw

When testing code that uses `libplumraw`, create a `TestWebConnection{}` object
instead of using `NewWebConnection{}` and hand that to your code to use as the
`WebConnection`. You can populate the `TestWebConnection{}` with House, Room,
LogicalLoad, etc. objects and when calling its methods it will return those
objects instead of making calls out to the website. This will give you
predictable responses you can use to exercise your code.  For Lightpads, use a
`TestLightpad{}` struct instead of the actual one with similar effects.

*/
package libplumraw
