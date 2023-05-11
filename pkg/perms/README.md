# pkg/permissions

Ideas and some code snippets have been taken from the [Permify/go-role repository](https://github.com/Permify/go-role),
which is licensed under [MIT license](https://github.com/Permify/go-role/blob/fe5a762e0605e42a246368dee9c54d2b28723dd0/LICENSE).

## Rework

- [ ] Permissions should be typed in some way (custom string type(s))
- [ ] Inheritance in a job/faction
- [ ] Permissions can be `True`/`"Unset"`/`False`
  - [ ] `"Unset"` means that it be inherited from the next lower role
- [ ] "Default" access group for any job that doesn't have
- [ ] Per Role/ Rank "Props"
