# AmVillage

AmVillage is a mass trading game designed to be played by around 100 players
live.

This application facilitates the trading part of the game.

Admins can add/remove items from arbitrary teams with the only limitation being
the balance not going under zero.

Players may trade with other teams based on the resources they have. However,
only one player can trade at a time to prevent situations where someone
successfully negotiates a trade but fails to do so as their teammate has
already traded that away.

Each player may request a lock for 60 seconds, after which they have to request
the lock again to trade. There are no mechanisms limiting the number of lock
acquiration as the players in the same team are assumed to be in contact with
each other.
