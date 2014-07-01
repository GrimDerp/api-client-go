api-client-go  |Build Status|_
==============================

.. |Build Status| image:: https://travis-ci.org/googlegenomics/api-client-go.png?branch=master
.. _Build Status: https://travis-ci.org/googlegenomics/api-client-go


The file ``main.go`` demonstrates how easy it is to use the `Google Genomics
API`_ with the `Go programming language`_.

.. _Google Genomics Api: https://developers.google.com/genomics/
.. _Go programming language: http://www.golang.org

Getting started
---------------

* First, you'll need a valid client ID and secret. Follow the `sign up
  instructions <https://developers.google.com/genomics>`_ and
  download the JSON file for ``installed application``.

* Install the example with::

   go get github.com/googlegenomics/api-client-go

* To see supported commands, run the program with::

   api-client-go

* To know about the flags in a particular command, run with ``help``, e.g.::

   api-client-go help
   api-client-go help readsets
   api-client-go help readsets search

* To search for a readsets, run with::

   api-client-go readsets search --use-oauth=client_secret.json --dataset_ids \
     376902546192
