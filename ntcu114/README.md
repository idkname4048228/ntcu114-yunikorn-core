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
