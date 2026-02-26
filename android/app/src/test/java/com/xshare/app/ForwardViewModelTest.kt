package com.xshare.app

import com.xshare.corebridge.Bridge
import org.junit.Assert.assertEquals
import org.junit.Test

class ForwardViewModelTest {
    @Test
    fun startForward_updatesStateToRunning() {
        val vm = ForwardViewModel(FakeBridge())

        vm.startForward()

        assertEquals(ForwardViewModel.State.Running, vm.state.value)
    }

    private class FakeBridge : Bridge {
        override fun startForwarding(): Boolean = true
        override fun stopForwarding() = Unit
    }
}
