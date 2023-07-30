# scrooge-bot

Telegram bot for checking routine expenses - https://t.me/my_wife_is_scrooge_bot

## Set expense

The message should contain the category as a hashtag and amounts. Example: *5.55 10.20 #food 40 - the last one is in the supermarket* 
The message can contain many amounts but only one hashtag.

## Commands

- ```/calc``` - calculate expenses for current month
- ```/calc_week``` - calculate expenses for 7 days
- ```/calc_month``` - calculate expenses for 30 days
- ```/calc_quarter``` - calculate expenses for 90 days

Command result example:
```
coffee: 1042
food: 42
pet: 24
```
