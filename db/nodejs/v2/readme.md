https://chatgpt.com/c/679e5e7d-cb4c-8006-9eb7-50fb11db11ed

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
- 🟩 run compaction by endpoint and add test to chec that after copmaction we can read data
- 🟩 do the same for level0 + level 1 -> level1+level2 -> level2+level3
- 🟩 do new more simple alg for others levels 


🟨 index in memory should be by interval/range to minimize data in memeory

🟨 gitignore
🟨 typescript
🟨 entities map
🟨 e2e tests directly for DBServer (without webclient but via driver)