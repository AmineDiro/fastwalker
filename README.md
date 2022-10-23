# FastWalker

Concurrent file system walker written in GO.

### Benchmark

I tested against a serial recursive version written in python using `scandir()`.

The `scripts/benchmark.sh` runs a filewalk using the python against the go version and clears the FS cache between each run.

On my system running : 
* CPU Ryzen AMD Ryzen 3600X with 12 thread and a pool of 24 goroutines
* NVME drive : Crucial P2 CT250P2SSD8 SSD with ext4 filesystem  

```bash
> ./scripts/benchmark.sh /home/amine

vm.drop_caches = 3
Finished crawling 1001472 in 20.29
vm.drop_caches = 3
2022/10/23 16:20:09 Walking [/home/amine]
2022/10/23 16:20:12 999814 files in 2.234110s

```

We achived a **~10x speedup** 💯.

### TODO : 
- [ ] Figure out the difference btw file count python vs golang
- [ ] Check for symlinks
- [ ] Add callback func to execute user defined function while walking the directory




