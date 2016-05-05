> 9. Write up your experience in a final paper, with the following requirements:
>
>   a. 12 point, Helvetica font, 0.5 inch margins all around, single spaced, standard character spacing. Full name and student ID at the top, PDF format.
>
>   b. Introduction, at least 1/2 page of text (10% of grade)
>       1. Brief summary of the problem you are solving and the features of the problem that are hard. For example, if you were trying to make a bot that plays chess you could describe how a good game strategy is to plan moves several steps ahead with contingencies for what the other player does. Does your bot have a model of the whole board, or is it only focused on the open squares directly adjacent to squares your bot is occupying. Make sure to state which game you are designing a bot for.

>    c. Methods, at least 1 or more pages of text, excluding pictures and code snippets. (50% of grade)
>
>       1. Detailed descriptions of the programming approaches you took, including those that failed. You should explore many ideas, and explain why you made the design decisions you did. For example, if you used an approach like depth first search, why did you think it would produce a successful bot (win games)?

>       2. Annotated code snippets should be included as examples of what you tried.

>       3. Note that exploratory or creative attempts using general AI methods (e.g., search, RL, NN, GA/GP, CSP, etc) are worth more much credit for the assignment than rank on the site, though both creativity and rank will be taken into consideration. I encourage you to test a variety of the AI methods we learned about this semester, and to include them in this section.

>    d. Results, 1/2 page of text, excluding pictures. (20% of paper grade (sub 5% based on rank on the site))

>       1. Report which of the approaches described in the methods section worked best and what measure made you conclude that it was the best.

>       2. Screenshot of your username and bot performing the game on aigames (which you reference in the paper)

>       3. Screenshot of your post in the discussion forum posting your gitlab repository (or wait until after really competing in blockbattle).

>       4. Links to your gitlab code repository

>    e. Discussion and conclusions, 1/2 to 2 pages as needed. (15% of paper grade)
>       1. Your remarks on performance of various strategies tried, ideas about the gameplay, future possibilities for improving your bot, and any other thoughts. What worked? What would you implement in the future if you were to try to keep improving the bot, or make it play another game?

---




<!--```
Policy: State -> Action
State: (update game, update player1, update player2)
Action: [Move] | None
Move: down | left | right | turn left | turn right | drop | skip
```-->

<!--There are a couple of ways to look at this problem of learning. The simplest is to consider look at playing a game as a stochastic function that takes in two policies and returns the one the that wins.
-->
<!--```
Game: (Policy 1, Policy 2) -> (1 | 2 | None)
```-->

<!--Policies can almost be looked at is a partially ordered with respect to the game function, but they are not, because the game function is stochastic. Still, we can look at this task as finding the policy most likely to dominate all other policies, i.e. the one that is the best.
-->
<!--Alternatively, we could break up each game into the separate turns. For each turn, we pass in our moves and the state, and get back the new state.
-->

<!--
```
Turn: (State, Action) -> State
```-->

<!--So if we have a way of estimating the value of extracting features from the state, we could use reinforcement learning to maximize the reward over time.
-->
<!--I stuck to the former, because it eliminated the need to approximate mid game rewards.
-->
<!--

To tried a couple of different methods to learn this policy function. I do not hand code any representation of the state, instead I let the learning methods figure out what sort of model is needed.  -->

<!-- The first, is to perform supervised learning to do offline learning of the policy. With this, we have a test policy, run a game using the test policy, and then update the test policy to reflect the results of the game. We could run two different policies against one another to get an idea of their relative performance.
-->
