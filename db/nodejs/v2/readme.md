https://chatgpt.com/c/679e5e7d-cb4c-8006-9eb7-50fb11db11ed

For run DB Server:
```bash
pnpm start:dev
```

For run DB Server:
```bash
pnpm start:debug
```




✅ red black tree
✅ memtable base
✅ samuraiDB base
✅ wrap db with server
🟩 flush: save data from memtable to file sstables 
- ✅ restore id-max index after restart
- ✅ restore ss tables indexes
- ✅ search in mem table and the in ss tables (from newwest to oldest (не факт))

✅ e2e webclient tests

🟩 compaction
- ✅ done compaction first value for first level 
- ✅ run compaction by endpoint and add test to chec that after copmaction we can read data
- ✅ delete item by key
- ✅ during compaction skip __DELETED__ itmas (now they are coping to level1 but dont must)
- ✅ delete all data for tests
- ✅ to e2e tests after compaction add new data and check that it's OK
- 🟩 do the same for level0 + level 1 -> level1+level2 -> level2+level3
- 🟩 run compaction by scheduler (пока без блокировок в последоватлеьном режиме)
- 🟨 unit tests???? for all components?
- 🟩 do new more simple alg for others levels (вроде нет смысла потому что по факту все слои всё равно нужно сливать с предыдущим.. а значит нужны теже курсоры...)

- 🟩 concurrency\multithreading

🟨 index in memory should be by interval/range to minimize data in memeory

🟨 gitignore
🟨 typescript
🟨 entities map
🟨 e2e tests directly for DBServer (without webclient but via driver)