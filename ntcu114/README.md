# NTCU 114級 畢業專題

## 上一個目標：試著利用外接的程式檔案

實作請參考[筆記](https://hackmd.io/@ntcu-k8s/S1tny9dG0)

ACS110131 寫完了，這裡是[commit 連結](https://github.com/YouKaiWu/ntcu114-yunikorn-core/commit/16225ab0b28464b65699f67b0daae6f77e84f60f)

git commit message 請參考[這篇文章](https://wadehuanglearning.blogspot.com/2019/05/commit-commit-commit-why-what-commit.html)

[GOA 的 colab](https://colab.research.google.com/drive/1lP3TBnIsf6nGDB-Ft_8SyCyWjQOEv4wq#scrollTo=B-e9lphkix5z)

### 遇到的小問題
將 colab 上 python 的 GOA 移植到 go 上面後，兩者的最終結果有些差距。
python 的可以穩定算出最佳解(目前以公平性來看)，而 go 是沒辦法算出最佳解的，甚至沒有一個相對穩定的解答。

以 user 要求 [1, 4] 及 [3, 1] 來看，當資源限制為 [9, 18] 時，都能計算出 [3, 2] 的答案。
但將資源限制提高到 [18, 28] 時，卻只有 python 可以穩定計算出 [6, 4] 的答案。

下面是 python 的視覺化
![pyGOA](https://hackmd.io/_uploads/SkGfy5AOC.gif)

下面是 go 的視覺化
![go_GOA](https://hackmd.io/_uploads/Hyu-k9CdA.gif)

之後找個時間解決這個問題。

## 現在目標：試著將 GOA 應用在 YuniKorn 排程
成功搞定

將 addUser, addNode 加在 partition 裡面
去 entrypoint 初始化 GOA
在 schedule 那邊執行

### 使用的 yaml 檔
```
apiVersion: batch/v1
kind: Job
metadata:
    name: pi
    namespace: testjob
spec:
    template:
        metadata:
            labels:
                applicationId: "app1"
            annotations:
                yunikorn.apache.org/user.info: "
                {
                    \"user\": \"user1\",
                    \"groups\": [
                        \"developers\",
                        \"devops\"
                    ]
                }"
        spec:
            schedulerName: yunikorn
            containers:
              - name: pi
                image: perl:5.34.0
                resources:
                    requests:
                        memory: "256Mi"
                        cpu: "0.5"
                    limits:
                        memory: "512Mi"
                        cpu: "1"
                command: ["perl",  "-Mbignum=bpi", "-wle", "print bpi(1000)"]
            restartPolicy: Never
```
以及
```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      schedulerName: yunikorn
      containers:
      - name: nginx
        image: nginx:1.14.2
        resources:
          requests:
              memory: "256Mi"
              cpu: "0.5"
          limits:
              memory: "512Mi"
              cpu: "1"
        ports:
        - containerPort: 80
```

### 遇到的小問題
應該要重構一下 GOA ，像是 userCount 就一堆地方都有

有嘗試在 entrypoint 那邊修改排程的間隔，但好像沒有用，要再研究看看要在哪裡改