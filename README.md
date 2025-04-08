# PowerClassic

## What is it?
An experimental actor based game minecraft server targetting the classic protocol sloppily whipped together over the course of 2 days.

## Why?
While working on a game server for a game I and some others are developing, we had originally started with an actor framework.
Eventually pitfalls were discovered and we eventually reverted to a traditional concurrent engine (still in-house & written in golang).
So I created PowerClassic to provide a testing ground with an established game client, to hopefully resolve those issues initially discovered.

## Features
* You can see other players move
* 2 Basic events with cancelling
* Sloppy
* Horrid control flow
* Crash on disconnect

## Planned Features:
* Fix inherit structural flaws causing some circular dependecies
* Clean up the messages
    * Probably something like this: `type Message struct {Exec: func()}` to provide a threadsafe way to run code on actors without defining every message seperately and flip flopping between packet messages. This is similar to what we have done in our existing game server.

## Want to contribute?
Open up a PR
