;; used commets to separate tasks
;; 10)

(defn variable [name] (fn [listV] (get listV name)))

(defn constant [val] (fn [_] val))

(defn opConstr [op]
      (fn [& args]
          (fn [params]
              (apply op (mapv #(% params) args)))))

(defn div [& args] (/ (double (first args)) (double (second args))))

(def add (opConstr +))
(def subtract (opConstr -))
(def multiply (opConstr *))
(def divide (opConstr div))
(def negate (opConstr -))

(def arcTan2 (opConstr Math/atan2))
(def arcTan (opConstr Math/atan))

(def operations
  {'+ add
   '- subtract
   '* multiply
   '/ divide
   'negate negate
   'atan arcTan
   'atan2 arcTan2})

;; 11)

(defn constructor
      [ctor prototype]
      (fn [& args] (apply ctor {:prototype prototype} args)))

(defn proto-get
      ([this key] (proto-get this key nil))
      ([this key default]
       (cond
         (contains? this key) (this key)
         (contains? this :prototype) (proto-get (this :prototype) key default)
         :else default)))

(defn proto-call
      [this key & args]
      (apply (proto-get this key) this args))

(defn method
      [key] (fn [this & args] (apply proto-call this key args)))

(def evaluate (method :evaluate))
(def toString (method :toString))
(def toStringPostfix (method :toStringPostfix))

(def OperationBase
  {:toString (fn [this]
                 (str "("
                      (proto-get this :symb)
                      " "
                      (clojure.string/join
                        " " (map #(proto-call % :toString)
                                 (:args this)))
                      ")"))

   :toStringPostfix (fn [this]
                        (str "("
                             (clojure.string/join
                               " " (map #(proto-call % :toStringPostfix)
                                        (:args this)))
                             " "
                             (proto-get this :symb)
                             ")"))


   :evaluate (fn [this params]
                 (apply (proto-get this :operationToApply)
                        (map #(proto-call % :evaluate params) (:args this))))})

(defn operationFormer [opSign operator]
      (constructor
        (fn [this & args] (assoc this :args args))
        (assoc OperationBase
               :symb opSign
               :operationToApply operator)))

(def Add (operationFormer '+ +))
(def Subtract (operationFormer '- -))
(def Multiply (operationFormer '* *))
(def Divide (operationFormer '/ div))
(def Negate (operationFormer 'negate -))

(def ArcTan2 (operationFormer 'atan2 Math/atan2))
(def ArcTan (operationFormer 'atan Math/atan))

(def Constant
  (constructor (fn [this val]
                   (assoc this :val val))
               {:toString (fn [this] (str (:val this)))
                :evaluate (fn [this _] (:val this))
                :toStringPostfix (fn [this] (str (:val this)))}))

(def Variable
  (constructor (fn [this name]
                   (assoc this :name name))
               {:toString (fn [this] (str (:name this)))
                :evaluate (fn [this params] (get params (str (Character/toLowerCase (last (:name this))))))
                :toStringPostfix (fn [this] (str (:name this)))}))

(def exprObjects
  {'+ Add
   '- Subtract
   '* Multiply
   '/ Divide
   'negate Negate
   'atan2 ArcTan2
   'atan ArcTan})

(defn parseItemBase [forNum forSymb itemsList]
      (
        fn [line]
           (let [input (read-string line)]
                (letfn [(parse [expr]
                               (cond
                                 (number? expr)
                                 (forNum expr)

                                 (symbol? expr)
                                 (forSymb (name expr))

                                 (list? expr)
                                 (apply (get itemsList (first expr)) (map parse (rest expr)))))]
                       (parse input)))
           ))

(def parseObject (parseItemBase Constant Variable exprObjects))

;; (part for 10*))
(def parseFunction (parseItemBase constant variable operations))


;; 12)

(defn -return [value tail] {:value value :tail tail})
(def -valid? boolean)
(def -value :value)
(def -tail :tail)

(def _empty (partial partial -return))

(defn _char [p]
      (fn [[c & cs]]
          (if (and c (p c))
            (-return c cs))))

(defn _map [f]
      (fn [result] (if (-valid? result)
                     (-return (f (-value result)) (-tail result)))))

(defn _combine [f a b]
      (fn [input]
          (let [ar ((force a) input)]
               (if (-valid? ar)
                 ((_map (partial f (-value ar)))
                  ((force b) (-tail ar)))))))

(defn _either [a b]
      (fn [input]
          (let [ar (a input)]
               (if (-valid? ar)
                 ar
                 ((force b) input)))))

(defn _parser [parser]
      (let [pp (_combine (fn [v _] v) parser (_char #{\u0000}))]
           (fn [input] (-value (pp (str input \u0000))))))

(defn +char [chars]
      (_char (set chars)))

(defn +map [f parser]
      (comp (_map f) parser))

(def +parser _parser)

(def +ignore
  (partial +map (constantly 'ignore)))

(defn- iconj [coll value]
       (if (= value 'ignore)
         coll
         (conj coll value)))

(defn +seq [& parsers]
      (reduce (partial _combine iconj) (_empty []) parsers))

(defn +seqf [f & parsers]
      (+map (partial apply f) (apply +seq parsers)))

(defn +seqn [n & parsers]
      (apply +seqf #(nth %& n) parsers))

(defn +or [parser & parsers]
      (reduce _either parser parsers))

(defn +opt [parser]
      (+or parser (_empty nil)))

(defn +star [parser]
      (letfn [(rec [] (+or (+seqf cons parser (delay (rec))) (_empty ())))]
             (rec)))

(defn +plus [parser]
      (+seqf cons parser (+star parser)))

(defn +str [parser]
      (+map (partial apply str) parser))
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;

(def bracket_types [["(" ")"] ["[" "]"] ["{" "}"] ["<" ">"]])

(def parseObjectPostfix
  (let [*all-chars (mapv char (range 0 128))
        *digit (+char (apply str (filter #(Character/isDigit %) *all-chars)))
        *space (+char (apply str (filter #(Character/isWhitespace %) *all-chars)))
        *ws (+ignore (+star *space))

        *number
        (+map Constant
              (+map read-string (+seqf str
                                       (+opt (+char "-"))
                                       (+seqf str
                                              (+str (+plus *digit))
                                              (+opt (+seqf str (+char ".") (+str (+star *digit))))))))

        *varName (+map Variable (+str (+plus (+char "xyzXYZ"))))

        *operation (apply +or
                          (for [[symb operation] exprObjects]
                               (+seqf (constantly operation)
                                      (apply +seq (map (comp +char str) (name symb))))))]

       (letfn [(+value []
                       (+or *number *varName (delay (+brackets))))

               (+brackets []
                          (apply +or
                                 (for [[open close] bracket_types]
                                      (+seqn 1
                                             (+char open)
                                             *ws
                                             (+map (fn [[args op]] (apply op args))
                                                   (+seqf vector
                                                          (+plus (+seqn 0 *ws (delay (+value)) *ws))
                                                          *operation))
                                             *ws
                                             (+char close)))))]

              (+parser (+seqn 0 *ws (+value) *ws)))))
