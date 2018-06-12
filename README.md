# bam-file-replacer
bam-file-replacer.. what it does:  Search for all files matching a wildcard pattern and replace their file contents with the file contents contained in the 'template-file'. Settings defined in a json config file.




You'll also need the bam.json file in the same directory this app runs from.   

```
{
  "Template": {
    "template-file": "bam.txt",
    "destination": "C:\\Users\\SkitchN\\Documents\\goprojects\\src\\templates\\testing",
    "glob": "*.mkts2"
   }
}
```

**template-file:** text file containing the text you want to use to replace the text in all files

**destination:** the root folder that will be scanned for files matching the 'glob' 

**glob:** wild card pattern of file types.

[download executable here](https://github.com/niski84/bam-file-replacer/releases)

[download config file here](https://github.com/niski84/bam-file-replacer/blob/master/bam.json)
