import os
from time import perf_counter

def scantree(dir: str):
    """Recursively yield DirEntry objects for given directory."""
    if os.access(dir,os.F_OK | os.R_OK):
        for entry in os.scandir(dir):
            if entry.is_dir(follow_symlinks=False):
                yield from scantree(entry.path)
            else:
                yield entry


if __name__ == "__main__":
    root_path = "/home/amine"
    s = perf_counter()
    entries = list(scantree(root_path))
    e = perf_counter()
    print(f"Finished crawling {len(entries)} in {e-s:.2f}")    
