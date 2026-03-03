package com.xshare.app

import android.os.Bundle
import androidx.appcompat.app.AppCompatActivity

class MainActivity : AppCompatActivity() {
    
    private lateinit var viewModel: ForwardViewModel
    
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        
        viewModel = ForwardViewModel(CoreBridgeImpl())
        
        setupUI()
    }
    
    private fun setupUI() {
        // UI setup will be implemented
    }
}

class CoreBridgeImpl : CoreBridge {
    init {
        System.loadLibrary("xshare_core")
    }
    
    external override fun startForwarding(): Result<Unit>
    external override fun stopForwarding(): Result<Unit>
    external override fun getStats(): Result<ForwardStats>
    external override fun isConnected(): Boolean
}
