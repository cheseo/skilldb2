((nil
  (eval
   (lambda ()
     (when (string= (or (file-name-extension (or buffer-file-name "")) "")
                    "tmpl")
       (sgml-mode))))))
