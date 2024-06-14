package com.example.teachtechai.view.topic

import android.util.Log
import android.view.LayoutInflater
import android.view.ViewGroup
import androidx.recyclerview.widget.DiffUtil
import androidx.recyclerview.widget.ListAdapter
import androidx.recyclerview.widget.RecyclerView
import com.example.teachtechai.R
import com.example.teachtechai.data.response.TopicItem
import com.example.teachtechai.databinding.DrawerItemBinding

class DrawerAdapter() : ListAdapter<TopicItem, DrawerAdapter.ViewHolder>(DIFF_CALLBACK) {
    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): ViewHolder {
        val binding = DrawerItemBinding.inflate(LayoutInflater.from(parent.context), parent, false)
        return ViewHolder(binding)
    }

    override fun onBindViewHolder(holder: ViewHolder, position: Int) {
        val topicItem = getItem(position)
        Log.d("DrawerAdapter", "onBindViewHolder - Position: $position, Item: $topicItem")
        holder.bind(topicItem)
    }

    inner class ViewHolder(private val binding: DrawerItemBinding) : RecyclerView.ViewHolder(binding.root) {
        fun bind(topicItem: TopicItem) {
            binding.apply {
                Log.d("DrawerAdapter", "Binding TopicItem: $topicItem")
                drawerText.text = topicItem.topic?: ""
                itemView.setOnClickListener {
                    // Handle click, for example:
                    // val intent = Intent(itemView.context, DetailTopicActivity::class.java).apply {
                    //     putExtra("topicId", topicItem.id)
                    // }
                    // itemView.context.startActivity(intent)
                }
            }
        }
    }

    companion object {
        private val DIFF_CALLBACK = object : DiffUtil.ItemCallback<TopicItem>() {
            override fun areItemsTheSame(oldItem: TopicItem, newItem: TopicItem): Boolean {
                return oldItem.id == newItem.id
            }

            override fun areContentsTheSame(oldItem: TopicItem, newItem: TopicItem): Boolean {
                return oldItem == newItem
            }
        }
    }
}
