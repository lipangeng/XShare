package com.xshare.app

import android.os.Bundle
import androidx.activity.ComponentActivity
import com.xshare.corebridge.CoreBridge

class MainActivity : ComponentActivity() {
    private val viewModel = ForwardViewModel(CoreBridge)

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
    }

    override fun onDestroy() {
        viewModel.stopForward()
        super.onDestroy()
    }
}
