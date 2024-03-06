#!/usr/bin/python3

import pod5 as p5
import csv
import os
import sys
import multiprocessing
import math

def Chunks(l,n):

    return [l[i:i+n] for i in range(0, len(l), n)]


def Writer(sli,pod5_reads):

    for s in sli:

        reader=p5.Reader(s[0])

        try:

            read = next(reader.reads(selection=[s[1]]))
            new_read=read.to_read()
            pod5_reads.append(new_read)

        except:

            print(s)

        #writer.add_read(new_read)
        reader.close()


def process(pass_dir, in_file, out_file,threads, skip_dir):

    #store paths
    allf=[]

    #store reads
    manager = multiprocessing.Manager()
    pod5_reads=manager.list()
    in_dir=pass_dir

    with open(in_file) as file_in:

        tsv_file = csv.reader(file_in, delimiter="\t")

        for line in tsv_file:

            in_dir=skip_dir if skip_dir != "None" and "skip" in line[0] else pass_dir
            pod5_path=os.path.join(in_dir, line[0])
            allf.append((pod5_path,line[1]))

    chunk_size=len(allf)/threads
    slices=Chunks(allf,math.ceil(chunk_size))
    processes=[]

    for _,sli in enumerate(slices):

        p=multiprocessing.Process(target=Writer, args=(sli,pod5_reads))
        p.start()
        processes.append(p)
        
    for p in processes:
        
        p.join()

    if len(pod5_reads) != 0:

        writer=p5.Writer(out_file)
        writer.add_reads(pod5_reads)
        writer.close()
    
    print("missing:", len(allf) - len(pod5_reads))

if __name__ == '__main__':

    pass_dir=sys.argv[1]
    in_file=sys.argv[2]
    out_file=sys.argv[3]
    threads=int(sys.argv[4])
    skip_dir=sys.argv[5]
    process(pass_dir, in_file, out_file,threads, skip_dir)
