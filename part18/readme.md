# 关于堆栈

    什么是堆/栈
    在这里并不打算详细介绍堆栈，仅简单介绍本文所需的基础知识。如下：
    
    堆（Heap）：一般来讲是人为手动进行管理，手动申请、分配、释放。
                一般所涉及的内存大小并不定，一般会存放较大的对象。另外其分配相对慢，涉及到的指令动作也相对多
    栈（Stack）：由编译器进行管理，自动申请、分配、释放。
                一般不会太大，我们常见的函数参数（不同平台允许存放的数量不同），局部变量等等都会存放在栈上
                
# 什么是逃逸分析

    在编译程序优化理论中，逃逸分析是一种确定指针动态范围的方法，简单来说就是分析在程序的哪些地方可以访问到该指针
    
    通俗地讲，逃逸分析就是确定一个变量要放堆上还是栈上，规则如下：
    
    是否有在其他地方（非局部）被引用。只要有可能被引用了，那么它一定分配到堆上。否则分配到栈上
    即使没有被外部引用，但对象过大，无法存放在栈区上。依然有可能分配到堆上
    对此你可以理解为，逃逸分析是编译器用于决定变量分配到堆上还是栈上的一种行为
    
# 在什么阶段确立逃逸

    在编译阶段确立逃逸，注意并不是在运行时

# 为什么需要逃逸

    这个问题我们可以反过来想，如果变量都分配到堆上了会出现什么事情？例如：

    垃圾回收（GC）的压力不断增大
    申请、分配、回收内存的系统开销增大（相对于栈）
    动态分配产生一定量的内存碎片
    其实总的来说，就是频繁申请、分配堆内存是有一定 “代价” 的。
    会影响应用程序运行的效率，间接影响到整体系统。
    因此 “按需分配” 最大限度的灵活利用资源，才是正确的治理之道。这就是为什么需要逃逸分析的原因，你觉得呢？
    
# 怎么确定是否逃逸
    
    第一，通过编译器命令，就可以看到详细的逃逸分析过程。而指令集 -gcflags 用于将标识参数传递给 Go 编译器，涉及如下：
    
    -m 会打印出逃逸分析的优化策略，实际上最多总共可以用 4 个 -m，但是信息量较大，一般用 1 个就可以了
    
    -l 会禁用函数内联，在这里禁用掉 inline 能更好的观察逃逸情况，减少干扰
    
    $ go build -gcflags '-m -l' main.go
    第二，通过反编译命令查看
    
    $ go tool compile -S main.go
    注：可以通过 go tool compile -help 查看所有允许传递给编译器的标识参数
