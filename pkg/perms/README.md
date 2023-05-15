# pkg/perms

Ideas and some code snippets have been taken from the [Permify/go-role repository](https://github.com/Permify/go-role),
which is licensed under [MIT license](https://github.com/Permify/go-role/blob/fe5a762e0605e42a246368dee9c54d2b28723dd0/LICENSE).

## Rework

- [ ] Permissions can be `True`/`"Unset"`/`False`
  - [ ] `"Unset"` means that it be inherited from the next lower role
- [ ] Permissions should be typed in some way (custom string type(s)) at least in the backend
- [ ] Inheritance from lower ranking job
- [x] Per Role/ Rank "Attributes" for storing, e.g., list of jobs + ranks.
- [ ] "Default" access group for any job that doesn't have
