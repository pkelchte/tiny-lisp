;; N クイーン問題
;; 猪股・益崎「Schemeによる記号処理入門」森北出版 1994年
;; §8.2 の Scheme コードから

(defun node-expand (n lst)
  (if (= n 0) nil
    (cons (cons n lst)
          (node-expand (- n 1) lst))))

(defun safe? (lst)
  (let ((new (car lst))
        (hlst (cdr lst)))
    (if (null hlst) t
      (safe-aux? new (+ new 1) (- new 1) hlst))))

(defun safe-aux? (new up down hlst)
  (if (null hlst) t
    (let ((pos (car hlst)))
      (and (not (= pos new))
           (not (= pos up))
           (not (= pos down))
           (safe-aux? new (+ up 1) (- down 1) (cdr hlst))))))

(defun goal? (x n) (= (length x) n))

(defun nqueens (n)
  (let (lst solution x pop push search)
    (setq lst (node-expand n nil))
    (setq solution nil)
    (setq x nil)
    (defun pop () (let ((y (car lst)))
                    (setq lst (cdr lst))
                    y))
    (defun push (y) (setq lst (append y lst)))
    (defun search ()
      (if (null lst)
          solution
        (setq x (pop))
;;      (princ "lst = ") (princ lst)
;;      (princ " x = ") (princ x) (terpri)
        (if (safe? x)
            (if (goal? x n)
                (setq solution (cons x solution))
              (push (node-expand n x))))
        (search)))
    (search)))

(print (length (nqueens 8)))            ; 92
