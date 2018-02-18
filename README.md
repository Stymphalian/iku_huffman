## Go Library - Huffman Encoding
Date: 2018-02-17 \
Author: jordanyu

### codec/
Simple encoder and decoder classes allowing you to write a 'payload' of ASCII
to any io.Writer stream

### huffman/
Contains the real code for reading and writing a huffman encoded payload

* **reader.go** - Contains the Reader class for reading a binary encoded payload 

* **writer.go** - Contains the Writer class for writing an ASCII payload into an huffman encoded form.

* **model.go** -  Model is the main structure which the huffman tree as well as a map from ascii symbols to their huffman bit patterns.

* **codebook.go** - Contains struct and functions used to represent the huffman codebook. This includes this like the frequency dictionary as well as the in-memory implementation of the huffman tree. There exists also methods for creating the canonical form of the huffman codebook.