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

* First, you'll need a valid client ID and secret. Follow the `authentication
  instructions <https://developers.google.com/genomics#authenticate>`_ and
  download the JSON file for ``installed application``.

* Install the example with::

   go get github.com/googlegenomics/api-client-go

* Run the program with::

   api-client-go -use_oauth client_secret.json


