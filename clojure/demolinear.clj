(defn unFab [f] (fn [& args] (apply mapv f args)))

(def v+ (unFab +))
(def v- (unFab -))
(def v* (unFab *))
(def vd (unFab /))

(defn scalar [& args] (reduce + (apply v* args)))

(defn vect [& args]
      (let [
            [a b] args
            [a1 a2 a3] a
            [b1 b2 b3] b
            ]
           [
            (- (* a2 b3) (* a3 b2))
            (- (* a3 b1) (* a1 b3))
            (- (* a1 b2) (* a2 b1))
            ]
           )
      )

(def m+ (unFab v+))
(def m- (unFab v-))
(def m* (unFab v*))
(def md (unFab vd))

(defn lowerMapving [f]
      (fn [first second]
          (mapv #(f % second) first)))

(def v*s (lowerMapving *))
(def m*s (lowerMapving v*s))
(def m*v (lowerMapving scalar))

(defn transpose [m] (apply mapv vector m))

(defn m*m [m1 m2]
      (let [m2_t (transpose m2)]
           (mapv (fn [row] (mapv (fn [col]
                                     (scalar row col))
                                 m2_t))
                 m1)))

(defn ordinaryVector [v]
      (and (vector? v) (every? number? v)))

(defn ordinaryMatrix [m]
      (and (vector? m) (every? ordinaryVector m)))

(defn tFab [f1 f2 f3]
      ( fn [& args]
           (cond
             (every? number? args)
             (apply f1 args)
             (every? ordinaryVector args)
             (apply f2 args)
             (every? ordinaryMatrix args)
             (apply f3 args)
             :else
             (apply mapv (tFab f1 f2 f3) args)))
      )

(def t+ (tFab + v+ m+))
(def t- (tFab - v- m-))
(def t* (tFab * v* m*))
(def td (tFab / vd md))
