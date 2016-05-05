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

Saul Shanabrook

29105706

# Block Battle Bot

## Introduction

I am designing a bot for the Block Battle challenge on The AI Games. The game is a competitive dual version of Tetris. The goal is to live longer than your opponent. If you perform certain special moves (like clearing multiple rows at once), you earn points which cause the opponent's board to fill up with added lines. The game must end, because the bottom rows are slowly filled up by the engine over time.

This game has a few challenges to it. First, it is not obvious how well you are doing, until the game is over. Although there is a "score", it does not determine the winner. It would not be difficult to win the game without getting any points, by just clearing rows individually. Also, since you are playing against another bot, you can't score your bot absolutely, only in relation to the other player. It is also stochastic, because the engine randomly chooses a block to add each round. There might also be a trade off of defense and offense, as your bought weighs clearing its own board versus gaining points to hinder the opponent.

If we look at every board and new piece combination as a possible state, the number of states is very large. It has an upper bound of roughly 2^(10 * 20) (each block in the board could be filled or not) * 5 (for every new piece) * 10 (for possible starting positions of the new pieces). This means we need some way of reducing the dimensionality of the state space, in order to store information about states in a generalizable manner.

## Methods

I chose to design to create a bot that would execute a learned policy function, which returns an action given a state. The state contains the current game information (fields for both players and the current piece and location). The action is the list of moves to take, to get the current piece into the desired location on the board.


I chose to use reinforcement learning to train the policy function. I was inspired by DeepMind's use of deep Q networks to play Atari games (insert cite). I quickly realized that it would make sense to use neural networks as well as the value function, because then I would not have to hand engineer the features. Neural networks have been used with great success in visual feature extraction and the Tetris state is similar to a low dimensional photograph. DeepMind had success with this in Atari games, which are also very similar to Tetris.

My approach has been to replicate their algorithm as much as possible, unless doing so wouldn't make sense for my problem or the tools I am using. I chose to program the learning and bot in the Go language, because of it's static typing and concurrency support. Most of my time was spent engineering the APIs that wrapped the game engine. I designed these in a way to support easily running many games concurrently and training on both players in each game. All games/players share the same neural network weights and experience bank. I make heavy use of Go's channels, and a pipelines model (insert cit), to keep this code clean and safe. Below you can see how the code for the core reinforcement learning run:

<insert code for running Q learning>

I chose to implement the "minibatch training"  differently than the DeepMind paper. The idea is that instead of updating the weights directly after having an experience (state, action, next state, reward), you instead save that experience in memory in a bounded length set. Then you sample randomly from the set and train the weights on those. They did those for a couple of reasons. First they said it reduces bias if we only end up in subset of states. Also, it allows use to learn many times from one experience. Instead of learning from the experience set every `n` rounds, I just continuously learned in the background, picking off an experience and learning from it. Since both the set of experiences and the neural networks weights were safe for concurrent read and writes (managed with channels and coroutines), doing this sped up the runs by a factor of 3, because it didn't have to wait to train to send out new actions.

Overall I managed to get the per game during training down to <insert seconds>. Below you can see the times for running different types of games, to see where the most time is spent:

<insert benchmark code>

I planned to use a methauristic search technique like genetic algorithms to optimize the metaparameters for the reinforcement learning, but I did not have time to do this.  I also wanted to train to use genetic algorithms to evolve the weights in the neural networks, as opposed to reinforcement learning, and see if that worked well.



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
